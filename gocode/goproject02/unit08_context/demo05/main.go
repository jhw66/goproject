//go:build go1.20
// +build go1.20

// demo05：WithCancelCause / WithTimeoutCause 与 context.Cause（需 Go 1.20+）。
// 取消时可附带原因，便于日志与排查链路与 ctx.Err() 的关系。
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	// WithTimeoutCause：超时原因可被 Cause 读取。
	ctx1, cancel1 := context.WithTimeoutCause(context.Background(), 100*time.Millisecond, errors.New("budget exhausted"))
	defer cancel1()
	<-ctx1.Done()
	fmt.Println("WithTimeoutCause Err():", ctx1.Err())
	fmt.Println("Cause:", context.Cause(ctx1))

	// WithCancelCause：手动取消并写入原因。
	ctx2, cancel2 := context.WithCancelCause(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel2(errors.New("upstream closed"))
	}()
	<-ctx2.Done()
	fmt.Println("WithCancelCause Err():", ctx2.Err())
	fmt.Println("Cause:", context.Cause(ctx2))
}
