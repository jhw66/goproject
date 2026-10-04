// WebSocket 广播服务：任意客户端发来的文本帧会转发给当前所有已连接客户端。
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	listenAddr = ":8888"
	wsPath     = "/ws"
)

// // Hub 在单 goroutine 内管理连接与广播，避免对 clients 的锁竞争。
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsMsg struct {
	mt   int
	data []byte
}

type Hub struct {
	clients    map[*Client]struct{}
	broadcast  chan wsMsg
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan wsMsg, 256),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
	}
}

// Client 每个连接一对读写 goroutine；只把读到的消息交给 Hub 广播。
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan wsMsg
}

const sendBuf = 256

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- wsMsg{mt: mt, data: message}
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(msg.mt, msg.data); err != nil {
			return
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case wsMsg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- wsMsg:
				default:
					// 发送队列满则跳过该客户端，避免阻塞整个 Hub
				}
			}
		}
	}
}

func serveWS(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade: %v", err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan wsMsg, sendBuf)}
	hub.register <- client
	go client.writePump()
	go client.readPump()
}

func main() {
	hub := newHub()
	go hub.run()

	r := gin.Default()
	r.GET(wsPath, func(c *gin.Context) { serveWS(hub, c) })
	log.Printf("broadcast server listening on %s  GET %s", listenAddr, wsPath)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
