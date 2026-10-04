// Hub 在单 goroutine 内管理连接与广播；来自 Redis Pub/Sub 与本地逻辑统一写入 broadcast。
package main

import (
	"github.com/gorilla/websocket"
)

type wsMsg struct {
	mt   int
	data []byte
}

// Hub 管理所有 WebSocket 客户端；broadcast 同时接收 Redis 订阅消息与（可选）本地注入。
type Hub struct {
	clients    map[*Client]struct{}
	broadcast  chan wsMsg
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan wsMsg, 512),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
				}
			}
		}
	}
}

// PublishText 将文本消息投递给 Hub，最终广播给所有已连接客户端（通常由 Redis 订阅回调调用）。
func (h *Hub) PublishText(data []byte) {
	h.broadcast <- wsMsg{mt: websocket.TextMessage, data: data}
}
