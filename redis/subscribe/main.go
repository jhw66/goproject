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
	sub := client.Subscribe(ctx, "channel_1")
	for ch := range sub.Channel() {
		fmt.Println(ch.Channel)
		fmt.Println(ch.Payload)
	}

	// for {
	// 	message, err := sub.ReceiveMessage(ctx)
	// 	if err == nil {
	// 		fmt.Println(message.Channel)
	// 		fmt.Println(message.Payload)
	// 	}
	// }

	// sub = client.PSubscribe(ctx, "channel_*")

	// for {
	// 	message, err := sub.ReceiveMessage(ctx)
	// 	if err == nil {
	// 		fmt.Println(message.Channel)
	// 		fmt.Println(message.Payload)
	// 	}
	// }

	sub.Unsubscribe(ctx, "channel_1")
	//sub.PUnsubscribe(ctx, "channel_*")
}
