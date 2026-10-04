package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// send buffer size
	bufSize = 256
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	roomID   string
	userID   string
	nickname string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send      chan []byte
	closeOnce sync.Once
}

func (c *Client) closeSend() {
	c.closeOnce.Do(func() {
		close(c.send)
	})
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		select {
		case c.hub.unregister <- c:
		case <-c.hub.done:
		}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		var in IncomingMessage
		err := c.conn.ReadJSON(&in)
		if err != nil {
			if !errors.Is(err, websocket.ErrCloseSent) &&
				websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		if in.Type != MessageTypeChat {
			continue
		}
		content := strings.TrimSpace(in.Content)
		if content == "" {
			continue
		}
		c.hub.Inbound(ChatMessage{
			RoomID:    c.roomID,
			UserID:    c.userID,
			Nickname:  c.nickname,
			Content:   content,
			CreatedAt: time.Now().UTC(),
		})
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				w.Close()
				return
			}
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write(newline); err != nil {
					w.Close()
					return
				}
				if _, err := w.Write(<-c.send); err != nil {
					w.Close()
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, store *Store, cfg AppConfig, w http.ResponseWriter, r *http.Request) {
	roomID := strings.TrimSpace(r.URL.Query().Get("roomId"))
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if roomID == "" || token == "" {
		http.Error(w, "roomId and token are required", http.StatusBadRequest)
		return
	}

	claims, err := ParseUserClaims(token, cfg.JWTSecret)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if err := store.EnsureRoomAccess(r.Context(), roomID, claims.UserID); err != nil {
		if errors.Is(err, ErrForbiddenRoom) {
			http.Error(w, "forbidden room", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:      hub,
		roomID:   roomID,
		userID:   claims.UserID,
		nickname: claims.Nickname,
		conn:     conn,
		send:     make(chan []byte, bufSize),
	}
	select {
	case client.hub.register <- client:
	case <-client.hub.done:
		conn.Close()
		return
	}

	history, err := store.ListMessages(context.Background(), roomID, 50)
	if err == nil {
		for _, msg := range history {
			payload, mErr := json.Marshal(msg)
			if mErr != nil {
				continue
			}
			select {
			case client.send <- payload:
			case <-client.hub.done:
				conn.Close()
				return
			}
		}
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
