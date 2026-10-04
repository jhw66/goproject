package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 1. 创建任务队列
	tasks := make(chan int, 100)

	// 2. 创建固定数量的worker
	workerCount := 3
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg)
	}

	// 3. 向任务队列发送任务
	for i := 0; i < 10; i++ {
		tasks <- i
	}

	// 4. 关闭任务队列（表示没有更多任务）
	close(tasks)

	// 5. 等待所有worker完成任务
	wg.Wait()
	fmt.Println("所有任务完成")
}

// worker函数
func worker(id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		fmt.Printf("Worker %d 开始处理任务 %d\n", id, task)
		time.Sleep(time.Second) // 模拟耗时任务
		fmt.Printf("Worker %d 完成任务 %d\n", id, task)
	}

	fmt.Printf("Worker %d 退出\n", id)
}
