// package main

// import (
// 	"fmt"
// 	"sync"
// )

// var wg sync.WaitGroup = sync.WaitGroup{}
// var lock sync.RWMutex

// func writeData(intChan chan int) {
// 	defer wg.Done()
// 	lock.RLock()
// 	for i := 1; i <= 10; i++ {
// 		intChan <- i
// 		fmt.Println("写入的数据为：", i)
// 	}
// 	lock.RUnlock()
// 	close(intChan)
// }
// func readData(intChan chan int) {
// 	defer wg.Done()
// 	for v := range intChan {
// 		fmt.Println("读取的数据为：", v)
// 	}
// }

// func main() {
// 	intChan := make(chan int, 100)

// 	wg.Add(2)
// 	go writeData(intChan)
// 	go readData(intChan)
// 	fmt.Printf("hhh!\n")

// 	wg.Wait()
// }

package main

import (
	"fmt"
)

func printer(c chan int) {
	// 开始无限循环等待数据
	for {
		// 从channel中获取一个数据
		data := <-c
		// 将0视为数据结束
		if data == 0 {
			fmt.Println("ok")
			break
		}
		// 打印数据
		fmt.Println(data)
	}
	// 通知main已经结束循环(我搞定了!)
	c <- 0
}
func main() {
	// 创建一个channel
	c := make(chan int)
	// 并发执行printer, 传入channel
	go printer(c)
	for i := 1; i <= 10; i++ {
		// 将数据通过channel投送给printer
		c <- i
		<-c
	}
	// 通知并发的printer结束循环(没数据啦!)
	c <- 0
	// 等待printer结束(搞定喊我!)
	<-c
}
