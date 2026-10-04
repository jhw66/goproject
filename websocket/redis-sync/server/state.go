// Redis 中维护连接元数据：每个连接一个 Hash，并在集合 ws:connections 中登记 ID，便于查询在线规模。
package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	connKeyPrefix = "ws:conn:"
	connSetKey    = "ws:connections"
	connTTL       = 10 * time.Minute
)

// ConnMeta 写入 Redis 的字段，可按业务扩展（用户、房间等）。
type ConnMeta struct {
	User       string
	RemoteAddr string
}

type ConnRegistry struct {
	rdb *redis.Client
}

func NewConnRegistry(rdb *redis.Client) *ConnRegistry {
	return &ConnRegistry{rdb: rdb}
}

func connHashKey(connID string) string {
	return connKeyPrefix + connID
}

// Register 登记连接并设置过期时间；同一 connID 应对应唯一 WebSocket。
func (r *ConnRegistry) Register(ctx context.Context, connID string, meta ConnMeta) error {
	key := connHashKey(connID)
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key,
		"user", meta.User,
		"remote_addr", meta.RemoteAddr,
		"connected_at", strconv.FormatInt(time.Now().Unix(), 10),
		"last_seen", strconv.FormatInt(time.Now().Unix(), 10),
	)
	pipe.SAdd(ctx, connSetKey, connID)
	pipe.Expire(ctx, key, connTTL)
	pipe.Expire(ctx, connSetKey, connTTL)
	_, err := pipe.Exec(ctx)
	return err
}

// Touch 刷新连接 TTL 与 last_seen，在读消息或心跳时调用。
func (r *ConnRegistry) Touch(ctx context.Context, connID string) error {
	key := connHashKey(connID)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, "last_seen", now)
	pipe.Expire(ctx, key, connTTL)
	pipe.Expire(ctx, connSetKey, connTTL)
	_, err := pipe.Exec(ctx)
	return err
}

// Unregister 删除连接状态并从集合移除。
func (r *ConnRegistry) Unregister(ctx context.Context, connID string) error {
	key := connHashKey(connID)
	pipe := r.rdb.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, connSetKey, connID)
	_, err := pipe.Exec(ctx)
	return err
}

// OnlineCount 返回当前在线连接数，并清理集合中 Hash 已过期的连接 ID。
func (r *ConnRegistry) OnlineCount(ctx context.Context) (int64, error) {
	if err := r.PruneStale(ctx); err != nil {
		return 0, err
	}
	n, err := r.rdb.SCard(ctx, connSetKey).Result()
	if err != nil {
		return 0, fmt.Errorf("scard %s:%w", connSetKey, err)
	}
	return n, nil
}

// PruneStale 移除集合中已经没有对应 Hash 的连接 ID。
func (r *ConnRegistry) PruneStale(ctx context.Context) error {
	ids, err := r.rdb.SMembers(ctx, connSetKey).Result()
	if err != nil {
		return fmt.Errorf("semeber %s:%w", connSetKey, err)
	}
	if len(ids) == 0 {
		return nil
	}

	pipe := r.rdb.Pipeline()
	redisCmd := make(map[string]*redis.IntCmd, len(ids))
	for _, id := range ids {
		redisCmd[id] = pipe.Exists(ctx, connHashKey(id))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("exists conn hashes:%w", err)
	}

	stale := make([]interface{}, 0)
	for id, res := range redisCmd {
		if res.Val() == 0 {
			stale = append(stale, id)
		}
	}
	if len(stale) == 0 {
		return nil
	}

	if err := r.rdb.SRem(ctx, connSetKey, stale...).Err(); err != nil {
		return fmt.Errorf("srem stale connections:%w", err)
	}
	return nil
}
