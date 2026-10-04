// demo04：WithValue 传递请求级数据；使用非导出类型作键，避免与其他包键冲突。
// 勿用 WithValue 传递可选函数参数——仅用于随请求跨 API 边界的数据。
package main

import (
	"context"
	"fmt"
)

// 私有类型键，外部包无法意外使用相同字面量键。
type ctxKey int

const (
	keyRequestID ctxKey = iota
	keyUserName
)

func handle(ctx context.Context) {
	id, _ := ctx.Value(keyRequestID).(string)
	name, _ := ctx.Value(keyUserName).(string)
	fmt.Printf("handle: request_id=%q user=%q\n", id, name)
	inner(ctx)
}

func inner(ctx context.Context) {
	fmt.Println("inner 仍能读到:", ctx.Value(keyRequestID))
}

func main() {
	base := context.Background()
	ctx := context.WithValue(base, keyRequestID, "req-42")
	ctx = context.WithValue(ctx, keyUserName, "alice")
	handle(ctx)
}
