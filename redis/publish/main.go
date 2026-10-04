package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "",
	})
	client.Publish(ctx, "channel_1", "msg")
	client.Publish(ctx, "channel_2", "hello")
	chs, _ := client.PubSubNumSub(ctx, "channel_1").Result()
	for ch, count := range chs {
		fmt.Println(ch)
		fmt.Println(count)
	}
	defer client.Close()
}
