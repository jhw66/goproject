// demo01：Context 接口四件套 + 根上下文 Background / TODO。
// Context 接口：
//   Deadline() (time.Time, bool)  返回截止时间；无则 ok=false。
//   Done() <-chan struct{}        取消/超时时会被关闭的通道；不可取消的 ctx 返回 nil。
//   Err() error                   Done 关闭后给出原因：Canceled 或 DeadlineExceeded。
//   Value(key any) any            读取请求级数据。
package main

import (
	"context"
	"fmt"
	"time"
)

func inspect(name string, ctx context.Context) {
	deadline, ok := ctx.Deadline()
	fmt.Printf("[%s] Deadline ok=%v deadline=%v\n", name, ok, deadline)
	fmt.Printf("[%s] Done()==nil ? %v\n", name, ctx.Done() == nil)
	fmt.Printf("[%s] Err()=%v\n", name, ctx.Err())
	fmt.Printf("[%s] Value(\"anyKey\")=%v\n\n", name, ctx.Value("anyKey"))
}

func work(ctx context.Context, name string) {
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] 提前结束: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("[%s] 步骤 %d\n", name, i+1)
			time.Sleep(60 * time.Millisecond)
		}
	}
	fmt.Printf("[%s] 正常完成\n", name)
}

func main() {
	// Background：根上下文，永不取消、无截止时间、不携带值。常用于 main/init/测试根。
	bg := context.Background()
	inspect("Background", bg)

	// TODO：占位，语义表达“将来要换成真实请求里的 ctx”，不要用 nil 代替。
	todo := context.TODO()
	inspect("TODO", todo)

	// 演示：没有取消信号时，一个按 ctx.Done 工作的函数也能正常完成。
	work(bg, "Background")
}
