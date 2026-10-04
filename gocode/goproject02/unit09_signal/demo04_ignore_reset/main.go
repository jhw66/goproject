// demo04：Ignore / Reset / Ignored —— 显式忽略信号或恢复系统默认行为。
//
//   - Ignore：收到信号时什么也不做；会撤销此前对这些信号的 Notify。
//   - Reset：撤销 Notify，恢复该信号的默认处理（与 Ignore 不同：默认可能是退出、core 等）。
//   - Ignored：查询某信号当前是否被忽略。
//
// 流程：先 Ignore(os.Interrupt)，此时 Ctrl+C 被忽略；Sleep 结束后 Reset，再等待一小段时间，
// 此时再按 Ctrl+C 会触发默认退出。若无人按键，约 12 秒后自动退出。
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("Ignored(os.Interrupt) 在 Ignore 之前:", signal.Ignored(os.Interrupt))

	signal.Ignore(os.Interrupt)
	fmt.Println("已 Ignore(os.Interrupt)，接下来 5 秒内按 Ctrl+C 应无反应")
	fmt.Println("Ignored(os.Interrupt) 现在:", signal.Ignored(os.Interrupt))

	time.Sleep(5 * time.Second)

	signal.Reset(os.Interrupt)
	fmt.Println("已 Reset(os.Interrupt)，默认处理已恢复；此时按 Ctrl+C 会结束进程")
	fmt.Println("Ignored(os.Interrupt) 现在:", signal.Ignored(os.Interrupt))
	fmt.Println("等待最多 7 秒…")

	<-time.After(7 * time.Second)
	fmt.Println("时间到，退出")
}
