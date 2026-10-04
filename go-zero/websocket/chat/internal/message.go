package internal

import "time"

const (
	MessageTypeChat = "message"
)

// IncomingMessage defines WebSocket client input shape.
type IncomingMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// ChatMessage is broadcast and persisted content.
type ChatMessage struct {
	ID        int64     `json:"id"`
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
