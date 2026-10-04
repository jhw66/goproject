// WebSocket + Redis 示例：
// - 实时同步：客户端发来的文本通过 Redis PUBLISH 到频道，本机 SUBSCRIBE 收到后由 Hub 广播给所有连接（多实例部署时每实例订阅同一频道即可横向扩展）。
// - 状态：每个连接在 Redis 中存 Hash（用户、地址、时间戳），并在集合 ws:connections 登记；断线时清理。
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	listenAddr   = ":8888"
	wsPath       = "/ws"
	redisChannel = "ws:broadcast"
	sendBuf      = 256
	writeWait    = 10 * time.Second
	pongWait     = 10 * time.Second
	pingPeriod   = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan wsMsg
	connID string
	reg    *ConnRegistry
	rdb    *redis.Client
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
		if c.connID != "" {
			if err := c.reg.Unregister(context.Background(), c.connID); err != nil {
				log.Printf("unregister redis %s: %v", c.connID, err)
			}
		}
	}()

	c.conn.SetReadLimit(4096)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		fmt.Println("收到pong")
		if err := c.reg.Touch(ctx, c.connID); err != nil {
			log.Printf("touch %s: %v", c.connID, err)
		}
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		mt, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read: %v", err)
			}
			return
		}
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err := c.reg.Touch(ctx, c.connID); err != nil {
			log.Printf("touch %s: %v", c.connID, err)
		}

		if mt != websocket.TextMessage {
			continue
		}
		// 统一走 Redis Pub/Sub，便于多实例间同步；本机订阅会收到自己发布的消息并广播。
		if err := c.rdb.Publish(ctx, redisChannel, string(p)).Err(); err != nil {
			log.Printf("publish: %v", err)
			return
		}
		fmt.Println("server publish message to redis")
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				deadline := time.Now().Add(writeWait)
				_ = c.conn.SetWriteDeadline(deadline)
				_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(msg.mt, msg.data); err != nil {
				return
			}
			fmt.Println("server write message to client")
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(writeWait)); err != nil {
				return
			}
			fmt.Println("ping")
		}
	}
}

func serveWS(hub *Hub, reg *ConnRegistry, rdb *redis.Client, baseCtx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}
		connID := uuid.NewString()
		user := strings.TrimSpace(c.Query("user"))
		meta := ConnMeta{
			User:       user,
			RemoteAddr: c.ClientIP(),
		}
		if err := reg.Register(baseCtx, connID, meta); err != nil {
			log.Printf("redis register: %v", err)
			_ = conn.Close()
			return
		}

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan wsMsg, sendBuf),
			connID: connID,
			reg:    reg,
			rdb:    rdb,
		}
		hub.register <- client
		go client.writePump()
		go client.readPump(baseCtx)
	}
}

func pubsubLoop(ctx context.Context, rdb *redis.Client, hub *Hub) {
	sub := rdb.Subscribe(ctx, redisChannel)
	defer sub.Close()

	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if msg == nil {
				continue
			}
			// 仅处理消息频道（忽略订阅确认等）
			if msg.Channel == redisChannel && msg.Payload != "" {
				hub.PublishText([]byte(msg.Payload))
			}
		}
	}
}

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr()})
	baseCtx := context.Background()
	if err := rdb.Ping(baseCtx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}

	hub := newHub()
	go hub.run()

	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()
	go pubsubLoop(ctx, rdb, hub)

	reg := NewConnRegistry(rdb)

	r := gin.Default()
	r.GET(wsPath, serveWS(hub, reg, rdb, baseCtx))
	r.GET("/stats/online", func(c *gin.Context) {
		n, err := reg.OnlineCount(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"online": n})
	})

	log.Printf("redis-sync server listen %s  GET %s  REDIS_ADDR=%s  PUB/SUB channel=%s",
		listenAddr, wsPath, redisAddr(), redisChannel)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	cancel()
}

func redisAddr() string {
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		return v
	}
	return "127.0.0.1:6379"
}
