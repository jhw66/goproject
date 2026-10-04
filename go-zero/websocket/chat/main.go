package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"websocket/chat/internal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

var (
	port    = flag.Int("port", 3333, "the port to listen")
	timeout = flag.Int64("timeout", 0, "timeout of milliseconds")
	cpu     = flag.Int64("cpu", 500, "cpu threshold")
)

func main() {
	flag.Parse()
	cfg := internal.LoadConfig(*port, *timeout, *cpu)

	logx.Disable()
	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Host:         cfg.Host,
		Port:         cfg.Port,
		Timeout:      cfg.Timeout,
		CpuThreshold: cfg.CpuThreshold,
	})
	defer engine.Stop()

	//连接mysql数据库
	db, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		panic(err)
	}

	//创建mysql存储实例
	store := internal.NewStore(db)

	//连接redis与redis存储实例
	broker := internal.NewRedisBroker(cfg)
	defer broker.Close()

	//创建hub实例
	hub := internal.NewHub(store, broker)
	defer hub.Shutdown()

	//启动hub
	go hub.Run()

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			if r.Method != "GET" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			http.ServeFile(w, r, "home.html")
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/ws",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			internal.ServeWs(hub, store, cfg, w, r)
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/rooms",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			claims, err := internal.AuthFromRequest(r, cfg)
			if err != nil {
				internal.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			rooms, err := store.ListVisibleRooms(r.Context(), claims.UserID)
			if err != nil {
				internal.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			internal.WriteJSON(w, http.StatusOK, map[string]interface{}{"rooms": rooms})
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/rooms",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			claims, err := internal.AuthFromRequest(r, cfg)
			if err != nil {
				internal.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			var req struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Visibility string `json:"visibility"`
				Password   string `json:"password"`
			}
			if err := internal.ReadJSON(r, &req); err != nil {
				internal.WriteError(w, http.StatusBadRequest, "invalid body")
				return
			}
			if strings.TrimSpace(req.ID) == "" {
				req.ID = internal.NewRoomID()
			}
			if strings.TrimSpace(req.Name) == "" {
				internal.WriteError(w, http.StatusBadRequest, "name required")
				return
			}
			if strings.TrimSpace(req.Visibility) == "" {
				req.Visibility = internal.VisibilityPublic
			}

			room := internal.Room{
				ID:         req.ID,
				Name:       req.Name,
				Visibility: req.Visibility,
				Password:   req.Password,
				CreatedBy:  claims.UserID,
			}
			if err := store.CreateRoom(r.Context(), room); err != nil {
				internal.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			if room.Visibility == internal.VisibilityPrivate {
				_ = store.AddMember(r.Context(), room.ID, claims.UserID, "owner")
			}
			internal.WriteJSON(w, http.StatusCreated, room)
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/rooms/:roomId/members",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			claims, err := internal.AuthFromRequest(r, cfg)
			if err != nil {
				internal.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			roomID := internal.ParseRoomID(r.URL.Path)
			if roomID == "" {
				internal.WriteError(w, http.StatusBadRequest, "roomId required")
				return
			}
			var req struct {
				UserID string `json:"userId"`
				Role   string `json:"role"`
			}
			if err := internal.ReadJSON(r, &req); err != nil {
				internal.WriteError(w, http.StatusBadRequest, "invalid body")
				return
			}
			if strings.TrimSpace(req.UserID) == "" {
				internal.WriteError(w, http.StatusBadRequest, "userId required")
				return
			}
			allowed, err := store.IsRoomOwner(r.Context(), roomID, claims.UserID)
			if err != nil {
				internal.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			if !allowed {
				internal.WriteError(w, http.StatusForbidden, "only room owner can add members")
				return
			}
			if err := store.AddMember(r.Context(), roomID, req.UserID, req.Role); err != nil {
				internal.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			internal.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodPost,
		Path:   "/rooms/:roomId/join",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			claims, err := internal.AuthFromRequest(r, cfg)
			if err != nil {
				internal.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			roomID := internal.ParseRoomID(r.URL.Path)
			if roomID == "" {
				internal.WriteError(w, http.StatusBadRequest, "roomId required")
				return
			}
			var req struct {
				Password string `json:"password"`
			}
			if err := internal.ReadJSON(r, &req); err != nil {
				internal.WriteError(w, http.StatusBadRequest, "invalid body")
				return
			}
			if err := store.JoinRoomByPassword(r.Context(), roomID, claims.UserID, req.Password); err != nil {
				if err == internal.ErrInvalidRoomPassword {
					internal.WriteError(w, http.StatusForbidden, "invalid room password")
					return
				}
				internal.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			internal.WriteJSON(w, http.StatusOK, map[string]string{"status": "joined"})
		},
	})

	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/rooms/:roomId/messages",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			claims, err := internal.AuthFromRequest(r, cfg)
			if err != nil {
				internal.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			roomID := internal.ParseRoomID(r.URL.Path)
			if roomID == "" {
				internal.WriteError(w, http.StatusBadRequest, "roomId required")
				return
			}
			if err := store.EnsureRoomAccess(r.Context(), roomID, claims.UserID); err != nil {
				internal.WriteError(w, http.StatusForbidden, "forbidden room")
				return
			}
			limit := internal.ParseLimit(r.URL.Query().Get("limit"), 50)
			msgs, err := store.ListMessages(r.Context(), roomID, limit)
			if err != nil {
				internal.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			internal.WriteJSON(w, http.StatusOK, map[string]interface{}{"messages": msgs})
		},
	})

	fmt.Printf("server started at http://%s:%d\n", cfg.Host, cfg.Port)
	engine.Start()
}
