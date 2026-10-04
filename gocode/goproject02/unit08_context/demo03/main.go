// demo03：WithTimeout / WithDeadline；超时后 ctx.Err() 为 DeadlineExceeded，
// 与 WithCancel 主动取消时的 Canceled 对比（包级变量）。
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	// WithTimeout：在指定时长后自动取消。
	ctx1, cancel1 := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel1()

	<-ctx1.Done()
	err1 := ctx1.Err()
	fmt.Println("WithTimeout 结束:", err1)
	fmt.Println("  errors.Is(..., DeadlineExceeded):", errors.Is(err1, context.DeadlineExceeded))
	fmt.Println("  errors.Is(..., Canceled):", errors.Is(err1, context.Canceled))

	// WithDeadline：在绝对时间点取消。
	deadline := time.Now().Add(200 * time.Millisecond)
	ctx2, cancel2 := context.WithDeadline(context.Background(), deadline)

	time.Sleep(100 * time.Millisecond)
	cancel2()
	<-ctx2.Done()
	err2 := ctx2.Err()
	fmt.Println("WithDeadline 结束:", err2)
	fmt.Println("  errors.Is(..., DeadlineExceeded):", errors.Is(err2, context.DeadlineExceeded))
	fmt.Println("  errors.Is(...,Canceled):", errors.Is(err2, context.Canceled))
}
