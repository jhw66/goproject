// 用锁来控制顺序，必须“双方都参与”

package main

import (
	"fmt"
	// "sync"
)

func main() {
	// lock := sync.Mutex{}
	// lock.Lock() // 主 goroutine 先上锁
	ch := make(chan int)
	c := make(chan int)

	go func(c chan int) {
		ch <- 1
		fmt.Println("Goroutine 收到数据并继续执行")
		// lock.Unlock() // 解锁给主 goroutine
		<-c
	}(c)

	<-ch
	// lock.Lock() // 等待子 goroutine 解锁
	fmt.Println("主 goroutine 发送数据后继续执行")
	// lock.Unlock()
	c <- 1
}

//-------------------------------------------------------------------

// 函数调用本身就是隐式的调度点

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	go func() {
// 		for i := 0; i < 50; i++ {
// 			fmt.Print("A") // 这里就是调度点！
// 		}
// 	}()
// 	go func() {
// 		for i := 0; i < 50; i++ {
// 			fmt.Print("B") // 这里也是调度点！
// 		}
// 	}()
// 	// 等待 goroutine 执行完成
// 	time.Sleep(time.Second)
// }

//-------------------------------------------------------------------
// 1. 循环是“顺序创建 goroutine”

// 	for i := 0; i < m; i++ {
// 	    go func(id int) { ... }(i)
// 	}

// 这段代码的执行主体仍然是主 goroutine。
// 主 goroutine 会顺序执行 i = 0, 1, 2 这三次循环。
// 每次执行到 go func(...) 时，会立即：
// 创建一个新的 goroutine；
// 把它加入 调度队列；
// 然后主 goroutine 继续执行下一次循环。
// 所以创建顺序是确定的：
// 👉 第 0 个 goroutine 先被创建，第 1 个其次，第 2 个最后。
// ⚙️ 2. 创建后，它们进入调度队列
// 当你调用 go f() 时，Go 并不会马上在当前线程执行 f()。
// 它会：
// 把这个 goroutine（即 f）放进 运行队列（run queue）；
// 等待 Go 调度器（runtime scheduler） 把它分配给某个 M（线程）执行。
// 这时它只是“被创建”，但是否马上执行、何时执行完全由调度器决定。

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	n := 10
// 	m := 3
// 	nums := make(chan int)
// 	for i := 0; i < m; i++ {
// 		go func(id int) {
// 			fmt.Println("create goroutine", id)
// 			for num := range nums {
// 				fmt.Printf("goroutine%d: %d\n", id, num)
// 			}
// 		}(i)
// 	}
// 	time.Sleep(time.Millisecond * 100) // 给 goroutine 启动时间
// 	for i := 1; i <= n; i++ {
// 		nums <- i
// 	}
// 	close(nums)
// 	time.Sleep(time.Second) // 等待打印完成
// }

// 解释
// 这种方式可以保证 数字顺序是正确的（1 到 n），
// 但打印的顺序可能被不同 goroutine 抢占（即输出顺序不一定稳定）。
// 所以这种方法不满足“打印顺序有序”，仅保证取数顺序有序。

//-------------------------------------------------------------------

// 数据竞争（data race）。

// package main

// func main() {
// 	x := 0
// 	for i := 0; i < 10; i++ {
// 		go func(x *int) {
// 			for i := 0; i < 10; i++ {
// 				*x++
// 			}
// 		}(&x)
// 	}
// 	// 等待子 goroutine 执行完成
// 	//time.Sleep(time.Second)
// 	println(x) // 这里打印的 x 不一定是 100
// 	// 因为主 goroutine 可能在子 goroutine 执行完成前就打印了 x
// 	// 也可能子 goroutine 执行了一部分后，主 goroutine 就打印了 x
// 	// 甚至可能子 goroutine 一条语句都没执行，主 goroutine 就打印了 x
// 	// 这完全取决于调度器的调度
// 	// 这种情况称为“数据竞争”（data race）
// }
