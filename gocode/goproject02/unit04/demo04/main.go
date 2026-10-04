package main

import (
	"fmt"
	"sync"
)

// 演示各种死锁场景
func main() {
	// 	主 goroutine 在等 <-ch（阻塞）
	// 但另一个 goroutine 在无限循环输出（没阻塞）
	// Go 运行时发现还有 goroutine 在运行 → 不算死锁 ✅

	// ch := make(chan int)
	// <-ch // 主 goroutine 阻塞
	// go func() {
	// 	for {
	// 		fmt.Println("我在运行中")
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	//-------------------------------------------------------

	// 	主 goroutine 阻塞在 <-ch；
	// 另一个 goroutine 运行着，并且会发送；
	// 当它执行 ch <- 42 时，会唤醒主 goroutine；
	// 所以最终能继续运行，不死锁 ✅

	// ch := make(chan int)
	// go func() {
	// 	ch <- 42 // 发送数据
	// }()
	// x := <-ch // 主 goroutine 阻塞等待
	// fmt.Println(x)

	//-------------------------------------------------------

	// 	主 goroutine 阻塞；
	// 其他 goroutine：没有；
	// 没人能唤醒主 goroutine；
	// 所以 Go 运行时检测到 “所有 goroutine 都 asleep”；
	// 于是报错 💥：
	// fatal error: all goroutines are asleep - deadlock!

	// ch := make(chan int)
	// x := <-ch // 没有任何人往 ch 里发送
	// fmt.Println(x)

	//-------------------------------------------------------

	// 	goroutine A 等待 <-ch2
	// goroutine B 等待 <-ch1
	// 两个都在等对方 → 循环等待
	// Go 检测到全体阻塞 → 死锁 💥。

	// ch1 := make(chan int)
	// ch2 := make(chan int)
	// go func() {
	// 	ch1 <- <-ch2
	// }()
	// ch2 <- <-ch1

	//-------------------------------------------------------

	// 	Go 调度器看到：
	// 主 goroutine：阻塞中
	// 其他 goroutine：还没创建（go func() 那行还没执行）
	// 因为主 goroutine 卡住，程序无法执行到下一行去启动接收的 goroutine。
	// 所以整个程序停在这里。没人能唤醒它。死锁。

	// ch := make(chan int)
	// ch <- 1 // (1) 发送，没有接收方
	// go func() {
	// 	<-ch
	// }()

	//-------------------------------------------------------

	// ┌─────────────┐         ┌──────────────┐
	// │ 子 goroutine │         │ 主 goroutine │
	// │  持有锁🔒     │         │  想拿锁       │
	// │  阻塞在 ch<- │         │  阻塞在 lock  │
	// │  等待主收数据│         │  等待子解锁   │
	// └──────┬──────┘         └──────┬──────┘
	//        │                       │
	//        └─────────死锁循环───────┘

	lock := sync.Mutex{}
	ch := make(chan int)
	go func() {
		lock.Lock()
		ch <- 1 // (1) 立即启动一个 goroutine，它阻塞在等待数据
		fmt.Println("Goroutine 收到数据并继续执行")
		lock.Unlock()
	}()
	lock.Lock()
	<-ch // (2) 主 goroutine 发送数据
	fmt.Println("主 goroutine 发送数据后继续执行")
	lock.Unlock()

	//-------------------------------------------------------

	// ch := make(chan int)

	// go func() {
	// 	for i := 0; i < 3; i++ {
	// 		ch <- i
	// 	}
	// 	// 如果这里不 close(ch)
	// 	// main 就会永远阻塞在 range
	// 	// close(ch)
	// }()

	// for v := range ch { // 永远等不到 close，循环不结束
	// 	fmt.Println(v)
	// }

}
