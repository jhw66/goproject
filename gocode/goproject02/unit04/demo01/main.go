// // 使用 channel 串联（信号传递法）
// // 这是最 Go 风格的写法。
// // 思路：用 m 个 goroutine + channel 串联，形成一个“接力棒”式打印机制。
// package main

// import (
// 	"fmt"
// )

// func main() {
// 	n := 10
// 	m := 3

// 	chans := make([]chan struct{}, m)
// 	//创建一个长度为 m 的切片，每个元素类型是 chan struct{}。
// 	// struct{} 是 空结构体，占用 0 字节内存，非常适合当信号用（类似于 C++ 的 semaphore / condition）。
// 	for i := range chans {
// 		chans[i] = make(chan struct{})
// 		//给每个位置分配一个 无缓冲通道。
// 		//无缓冲意味着：发送和接收必须同步发生，天然能实现“阻塞等待”的效果。
// 	}

// 	// 启动 m 个 worker
// 	for i := 0; i < m; i++ {
// 		go func(id int) {
// 			for num := id + 1; num <= n; num += m { //让每个 goroutine 打印它负责的那一组数字。
// 				<-chans[id]      //从自己的通道接收信号（阻塞等待）。只有当前一个 goroutine 向此通道发送信号后，才能继续。
// 				fmt.Println(num) // 打印
// 				chans[(id+1)%m] <- struct{}{}
// 				//通知下一个 goroutine（(id+1)%m 实现环形下标）。
// 				// 例如：
// 				// goroutine 0 通知 goroutine 1
// 				// goroutine 1 通知 goroutine 2
// 				// goroutine 2 通知 goroutine 0
// 			}
// 		}(i) //定义匿名函数并立即调用，go 关键字表示并发执行。
// 	}
// 	// 启动第一个 goroutine
// 	chans[0] <- struct{}{}
// 	// 空的 select {} 会永久阻塞主协程。
// 	// 这样主函数不会提前退出，让其他 goroutine 能一直执行。
// 	select {}
// }

// //struct{} —— 空结构体类型
// // 在 Go 中：
// // struct{}    // 表示一个“空结构体类型”，里面没有任何字段
// // 例如：
// // type Empty struct{} // 定义一个空结构体类型
// // ✅ 特点：占用内存为 0 字节。
// // Go 语言规范明确指出：一个空结构体在内存中不占空间。
// // 因此，它非常适合用来当作：
// // “信号”
// // “标志位”
// // “集合的键（set）中用作值类型”

// //struct{}{} —— 空结构体的一个值（实例）
// // 这其实是：
// // struct{}{} = 类型 + 字面值
// // 也就是说：
// // 前一个 {} 表示类型定义：空结构体类型 struct{}
// // 后一个 {} 表示该类型的字面量值
// // 举个例子：
// // var x struct{}   // 定义变量
// // x = struct{}{}   // 给它赋一个空结构体值
// // 因为空结构体没有字段，所以字面量也是空的 {}。

package main

import (
	"fmt"
	"sync"
)

func main() {
	n := 5
	m := 100
	wg := sync.WaitGroup{}
	wg.Add(n)

	chans := make([]chan struct{}, n)
	for i := range chans {
		chans[i] = make(chan struct{})
	}

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			for j := id; j <= m; j += n {
				<-chans[id]
				fmt.Println("当前goroutine为", id, "打印的数字为", j)

				// 如果下一个数字还存在，再通知下一个 goroutine
				if j+1 <= m {
					chans[(id+1)%n] <- struct{}{}
				}
			}
		}(i)
	}

	// 启动第一个 goroutine
	chans[0] <- struct{}{}

	wg.Wait()
}
