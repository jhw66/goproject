//go:build go1.21
// +build go1.21

// demo06：AfterFunc 在 ctx 取消后于独立 goroutine 中执行回调（需 Go 1.21+）。
// 仍应保证 WithXxx 返回的 cancel 被调用（通常 defer），避免泄漏子节点与定时器。
package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cleaned atomic.Bool
	stop1 := context.AfterFunc(ctx, func() {
		time.Sleep(50 * time.Millisecond)
		cleaned.Store(true)
		fmt.Println("AfterFunc: ctx 已取消，执行清理逻辑")
	})

	stop2 := context.AfterFunc(ctx, func() {
		time.Sleep(100 * time.Millisecond)
		cleaned.Store(false)
		fmt.Println("AfterFunc: ctx 已取消，执行清理逻辑")
	})
	defer stop2()

	<-ctx.Done()
	fmt.Println("清理已执行:", cleaned.Load())
	fmt.Println("主流程观察到:", ctx.Err())
	fmt.Println("stop1:", stop1())
	time.Sleep(50 * time.Millisecond)
	fmt.Println("清理已执行:", cleaned.Load())
	time.Sleep(50 * time.Millisecond)
	fmt.Println("清理已执行:", cleaned.Load())
}
