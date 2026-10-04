package cache

import (
	"context"
	"demo03/internal/config"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(c config.RedisConf) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConns,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis failed:%w", err)
	}

	return rdb, nil
}
