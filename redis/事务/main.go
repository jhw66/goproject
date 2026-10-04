package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	pipe := client.TxPipeline()
	incr := pipe.Incr(ctx, "tx_pipeline_counter")
	pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)
	_, err := pipe.Exec(ctx)
	fmt.Println(incr.Val(), err)
}
