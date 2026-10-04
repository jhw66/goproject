// JSON 消息示例：ReadJSON / WriteJSON、SetReadLimit、读超时、定时 Ping（WriteControl）、写锁。
package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ChatMsg 与客户端约定的 JSON 结构（字段可按业务扩展）。
type ChatMsg struct {
	Op   string `json:"op"`
	Text string `json:"text,omitempty"`
}

func handleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade:%v", err)
		return
	}
	defer conn.Close()

	conn.SetReadLimit(512 << 10)
	const pongWait = 20 * time.Second
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	var mu sync.Mutex
	writeJSON := func(v any) error {
		mu.Lock()
		defer mu.Unlock()
		return conn.WriteJSON(v)
	}
	writePing := func() error {
		mu.Lock()
		defer mu.Unlock()
		deadline := time.Now().Add(10 * time.Second)
		return conn.WriteControl(websocket.PingMessage, nil, deadline)
	}

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(18 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := writePing(); err != nil {
					return
				}
			}
		}
	}()

	for {
		var msg ChatMsg
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		conn.SetReadDeadline(time.Now().Add(pongWait))

		var werr error
		switch msg.Op {
		case "echo":
			werr = writeJSON(ChatMsg{Op: "echo", Text: msg.Text})
		default:
			werr = writeJSON(ChatMsg{Op: "error", Text: "unknown op"})
		}
		if werr != nil {
			log.Printf("writeJSON: %v", werr)
			break
		}
	}
	close(done)
}

func main() {
	r := gin.Default()
	r.GET("/ws", handleWS)
	log.Println("json_server listening on :8888  GET /ws")
	if err := r.Run(":8888"); err != nil {
		log.Fatal(err)
	}
}
