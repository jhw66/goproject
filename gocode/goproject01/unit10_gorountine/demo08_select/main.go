// package main

// import (
// 	"fmt"
// 	"time"
// )

// // 在 Go 语言中，select 语句的作用是：
// // 👉 在多个通道（channel）操作中同时等待，选择其中一个可以执行的分支去执行。
// // 🧠 一句话理解
// // select 是 通道版的 switch ——
// // 它让 goroutine 同时“监听”多个 channel，一旦某个 channel 可以收发数据，就执行对应的 case 分支。
// func main() {
// 	intChan := make(chan int, 1)
// 	go func() {
// 		fmt.Println("存int数据")
// 		intChan <- 10
// 	}()
// 	stringChan := make(chan string, 2)
// 	go func() {
// 		fmt.Println("存string数据")
// 		stringChan <- "asdfghjkl"
// 	}()

// 	// fmt.Println(<-intChan)
// 	select {
// 	case v := <-intChan:
// 		fmt.Println(v)
// 	case v := <-stringChan:
// 		fmt.Println(v)
// 	case stringChan <- "qwertyuiop":
// 		fmt.Println("存string数据2")
// 	default:
// 		fmt.Println("都阻塞了，执行default")
// 	}
// 	// 如果所有 case 都阻塞：
// 	// 有 default 分支时：执行 default；
// 	// 没有 default：select 会阻塞，直到有某个通道准备好

// 	ch1 := make(chan string)
// 	ch2 := make(chan string)

// 	go func() {
// 		time.Sleep(2 * time.Second)
// 		ch1 <- "来自 ch1 的消息"
// 	}()
// 	go func() {
// 		time.Sleep(1 * time.Second)
// 		ch2 <- "来自 ch2 的消息"
// 	}()

// 	select {
// 	case msg1 := <-ch1:
// 		fmt.Println(msg1)
// 	case msg2 := <-ch2:
// 		fmt.Println(msg2)
// 	case <-time.After(3 * time.Second):
// 		fmt.Println("超时退出")
// 	}
// }

// // 常见用途
// // 场景	作用
// // 多路通信	同时等待多个 goroutine 的输出
// // 超时控制	使用 time.After 实现超时等待
// // 停止信号	使用 done 通道优雅退出 goroutine
// // 非阻塞操作	使用 default 实现非阻塞发送或接收

package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	quit := make(chan bool)
	//新开一个协程
	go func() {
		for {
			select {
			case num := <-ch:
				fmt.Println("num = ", num)
			case <-time.After(3 * time.Second):
				fmt.Println("超时")
				quit <- true
			}
		}
	}()
	for i := 0; i < 5; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
	<-quit
	fmt.Println("程序结束")
}
