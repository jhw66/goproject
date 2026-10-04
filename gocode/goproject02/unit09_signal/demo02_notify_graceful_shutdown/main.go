// demo02：监听多种信号，做「优雅退出」的常见骨架。
//
//   - Unix：常见组合为 SIGINT + SIGTERM（例如 kill 默认发 TERM）。
//   - Windows：Notify(os.Interrupt) 处理 Ctrl+C；文档说明控制台关闭/注销/关机时也可能收到 syscall.SIGTERM。
//
// 运行后按 Ctrl+C，或等待约 8 秒由程序自行结束（演示清理逻辑）。
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	stop := make(chan struct{})
	go func() {
		tick := time.NewTicker(2 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-stop:
				return
			case t := <-tick.C:
				fmt.Println("[工作协程] 心跳", t.Format("15:04:05"))
			}
		}
	}()

	fmt.Println("模拟长期任务：按 Ctrl+C 或发送 SIGTERM 可触发优雅退出；否则 8 秒后自动结束。")

	select {
	case s := <-sigCh:
		fmt.Println("收到信号，开始收尾:", s)
	case <-time.After(8 * time.Second):
		fmt.Println("未到信号，演示时间到，开始收尾")
	}

	close(stop)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("清理完成，退出")
}
