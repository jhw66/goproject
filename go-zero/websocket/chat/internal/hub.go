package internal

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clientsByRoom map[string]map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	inbound  chan ChatMessage
	outbound chan ChatMessage
	done     chan struct{}
	stopped  chan struct{}

	broker *RedisBroker
	store  *Store

	mu     sync.Mutex
	cancel map[string]func()
	once   sync.Once
}

func NewHub(store *Store, broker *RedisBroker) *Hub {
	return &Hub{
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		inbound:       make(chan ChatMessage, 512),
		outbound:      make(chan ChatMessage, 512),
		done:          make(chan struct{}),
		stopped:       make(chan struct{}),
		clientsByRoom: make(map[string]map[*Client]bool),
		broker:        broker,
		store:         store,
		cancel:        make(map[string]func()),
	}
}

func (h *Hub) Run() {
	defer close(h.stopped)

	for {
		select {
		case client := <-h.register:
			roomClients := h.clientsByRoom[client.roomID]
			if roomClients == nil {
				roomClients = make(map[*Client]bool)
				h.clientsByRoom[client.roomID] = roomClients
				h.ensureRoomSubscription(client.roomID)
			}
			roomClients[client] = true
			h.broker.MarkOnline(context.Background(), client.roomID, client.userID)
		case client := <-h.unregister:
			h.removeClient(client)
		case message := <-h.inbound:
			saved, err := h.store.SaveMessage(context.Background(), message)
			if err != nil {
				log.Printf("save message failed: %v", err)
				continue
			}
			if err := h.broker.Publish(context.Background(), *saved); err != nil {
				log.Printf("publish message failed: %v", err)
			}
		case message := <-h.outbound:
			h.broadcast(message)
		case <-h.done:
			h.closeAllClients()
			h.cancelAllSubscriptions()
			return
		}
	}
}

func (h *Hub) ensureRoomSubscription(roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.cancel[roomID]; ok {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	stopSub := h.broker.SubscribeRoom(ctx, roomID, func(msg ChatMessage) {
		h.enqueueBroadcast(msg)
	})
	h.cancel[roomID] = func() {
		cancel()
		stopSub()
	}
}

func (h *Hub) cancelRoomSubscription(roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if cancel, ok := h.cancel[roomID]; ok {
		cancel()
		delete(h.cancel, roomID)
	}
}

func (h *Hub) enqueueBroadcast(msg ChatMessage) {
	select {
	case h.outbound <- msg:
	case <-h.done:
	default:
		log.Printf("outbound queue full for room=%s", msg.RoomID)
	}
}

func (h *Hub) broadcast(msg ChatMessage) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal outbound message failed: %v", err)
		return
	}

	roomClients := h.clientsByRoom[msg.RoomID]
	for client := range roomClients {
		h.broker.KeepOnline(context.Background(), client.roomID, client.userID)
		select {
		case client.send <- payload:
		default:
			h.removeClient(client)
		}
	}
}

func (h *Hub) Inbound(msg ChatMessage) {
	select {
	case h.inbound <- msg:
	case <-h.done:
	default:
		log.Printf("inbound queue full for room=%s", msg.RoomID)
	}
}

func (h *Hub) removeClient(client *Client) {
	if client == nil {
		return
	}
	roomClients, ok := h.clientsByRoom[client.roomID]
	if !ok {
		return
	}
	if _, exists := roomClients[client]; !exists {
		return
	}
	delete(roomClients, client)
	client.closeSend()
	h.broker.MarkOffline(context.Background(), client.roomID, client.userID)
	if len(roomClients) == 0 {
		delete(h.clientsByRoom, client.roomID)
		h.cancelRoomSubscription(client.roomID)
	}
}

func (h *Hub) closeAllClients() {
	for roomID, roomClients := range h.clientsByRoom {
		for client := range roomClients {
			client.closeSend()
			h.broker.MarkOffline(context.Background(), client.roomID, client.userID)
			delete(roomClients, client)
		}
		delete(h.clientsByRoom, roomID)
	}
}

func (h *Hub) cancelAllSubscriptions() {
	h.mu.Lock()
	cancellers := make([]func(), 0, len(h.cancel))
	for _, cancel := range h.cancel {
		cancellers = append(cancellers, cancel)
	}
	h.cancel = make(map[string]func())
	h.mu.Unlock()

	for _, cancel := range cancellers {
		cancel()
	}
}

func (h *Hub) Shutdown() {
	h.once.Do(func() {
		close(h.done)
		<-h.stopped
	})
}
