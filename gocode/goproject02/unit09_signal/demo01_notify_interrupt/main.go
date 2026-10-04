// demo01：signal.Notify 最基本用法 —— 把「异步信号」转发到 channel，替代进程的默认处理。
//
// 要点（见 https://pkg.go.dev/os/signal#Notify）：
//   - 必须用带缓冲的 channel（常见为 1），否则在尚未接收时可能丢通知。
//   - 未列出具体信号时，默认行为例如：SIGINT 会令程序退出；Notify 后改为投递到 c。
//   - Windows：对 os.Interrupt 使用 Notify 后，Ctrl+C / Ctrl+Break 会发到 channel，而不会直接退出。
//
// 运行本程序后按 Ctrl+C，应打印收到的信号；若 30 秒内未按键则超时退出（方便无人值守运行）。
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Println("已注册 os.Interrupt，请按 Ctrl+C（Windows 上 Ctrl+Break 也会触发 Interrupt）")
	fmt.Println("或等待 30 秒自动退出…")

	select {
	case s := <-c:
		fmt.Println("收到信号:", s)
	case <-time.After(30 * time.Second):
		fmt.Println("超时退出（未收到信号）")
	}
}
