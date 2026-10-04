// demo05：signal.NotifyContext（Go 1.16+）—— 把信号映射为可取消的 context。
//
// 要点（见 https://pkg.go.dev/os/signal#NotifyContext）：
//   - 返回的 ctx 在收到任一注册信号、调用 stop()、或 parent 取消时 Done。
//   - 必须在任务结束后调用 stop()，释放注册并可能恢复默认信号行为。
//   - Go 1.20+：若因信号取消，可用 context.Cause(ctx) 查看原因。
//
// 运行后可在 10 秒内按 Ctrl+C 观察 ctx 取消；否则超时打印后退出。
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func work(ctx context.Context) {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[work] 结束:", ctx.Err())
			if cause := context.Cause(ctx); cause != nil {
				fmt.Println("[work] context.Cause:", cause)
			}
			return
		case t := <-tick.C:
			fmt.Println("[work] tick", t.Format("15:04:05"))
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go work(ctx)

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("main: 10 秒未到信号，正常结束")
	case <-ctx.Done():
		fmt.Println("main:", ctx.Err())
		if cause := context.Cause(ctx); cause != nil {
			fmt.Println("main: context.Cause:", cause)
		}
		stop()
	}
	time.Sleep(200 * time.Millisecond)
}
