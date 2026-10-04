package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func roomChannel(roomID string) string {
	return fmt.Sprintf("chat:room:{%s}:channel", roomID)
}

func roomOnlineKey(roomID, userID string) string {
	return fmt.Sprintf("chat:room:{%s}:online:%s", roomID, userID)
}

// RedisBroker handles pub/sub and online presence.
type RedisBroker struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisBroker(cfg AppConfig) *RedisBroker {
	return &RedisBroker{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		}),
		ttl: time.Duration(cfg.OnlineTTLSeconds) * time.Second,
	}
}

func (b *RedisBroker) Close() error {
	return b.client.Close()
}

func (b *RedisBroker) Publish(ctx context.Context, msg ChatMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return b.client.Publish(ctx, roomChannel(msg.RoomID), payload).Err()
}

func (b *RedisBroker) MarkOnline(ctx context.Context, roomID, userID string) {
	if err := b.client.Set(ctx, roomOnlineKey(roomID, userID), "1", b.ttl).Err(); err != nil {
		log.Printf("mark online failed: %v", err)
	}
}

func (b *RedisBroker) MarkOffline(ctx context.Context, roomID, userID string) {
	if err := b.client.Del(ctx, roomOnlineKey(roomID, userID)).Err(); err != nil {
		log.Printf("mark offline failed: %v", err)
	}
}

func (b *RedisBroker) KeepOnline(ctx context.Context, roomID, userID string) {
	if err := b.client.Expire(ctx, roomOnlineKey(roomID, userID), b.ttl).Err(); err != nil {
		log.Printf("keep online failed: %v", err)
	}
}

func (b *RedisBroker) SubscribeRoom(ctx context.Context, roomID string, handler func(ChatMessage)) func() {
	sub := b.client.Subscribe(ctx, roomChannel(roomID))
	done := make(chan struct{})

	go func() {
		defer close(done)
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case payload, ok := <-ch:
				if !ok {
					return
				}
				var msg ChatMessage
				if err := json.Unmarshal([]byte(payload.Payload), &msg); err != nil {
					log.Printf("unmarshal room message failed: %v", err)
					continue
				}
				handler(msg)
			}
		}
	}()

	return func() {
		_ = sub.Close()
		<-done
	}
}
