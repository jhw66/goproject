// demo03：signal.Stop —— 取消此前对该 channel 的 Notify，恢复对应信号的默认处理。
//
// 文档：Stop 返回后，保证 c 不会再收到新的信号通知。
//
// 本例：先 Notify，再立刻 Stop，随后若你按 Ctrl+C，进程会按默认行为直接退出（本程序不会从 channel 读到信号）。
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

	signal.Stop(c)
	fmt.Println("已对 channel 调用 Stop：Ctrl+C 恢复为默认行为（进程会直接退出）。")
	fmt.Println("若 15 秒内不按 Ctrl+C，程序将正常结束。")

	select {
	case s := <-c:
		// Stop 之后通常走不到这里（信号不再投递到 c）
		fmt.Println("意外收到（不应出现）:", s)
	case <-time.After(15 * time.Second):
		fmt.Println("正常结束")
	}
}
