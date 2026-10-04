// 使用 sync.Mutex + Condition

// sync.Cond 是一个条件变量（Condition Variable），用来让 goroutine 在“某个条件未满足”时阻塞等待，当条件被满足后，再通过唤醒机制继续执行
// 1,cond.Wait():
// 当前 goroutine 会：
// 自动释放持有的锁；
// 阻塞等待其他 goroutine 调用 Signal() 或 Broadcast()；
// 被唤醒后，重新加锁；
// 返回到 Wait() 之后继续执行。
// ⚠️ 使用规则
// Wait() 必须在持有锁的前提下调用，否则会 panic。
// 2,cond.Signal()
// 唤醒一个正在 Wait() 的 goroutine（如果有）。
// 该被唤醒的 goroutine 会在 Signal() 返回后重新尝试获得锁；
// 如果当前锁仍被其他 goroutine 占用，它会继续阻塞；
// 因此 Signal() 只是发出“唤醒信号”，并非立即执行。
// 3,Broadcast()
// 唤醒所有正在等待的 goroutine。
// 通常用于：
// 状态变化影响所有等待者；
// 程序即将结束时唤醒所有阻塞协程退出。

//Cond 使用注意事项（重点）
// Wait() 必须在加锁状态下调用；
// 使用 for 循环检查条件；
// Wait() 内部会自动释放锁；
// 被唤醒后会重新获得锁；
// Signal() / Broadcast() 只是发出信号，不会立即解锁或强制执行；
// Cond 的锁一般是 *sync.Mutex，不要随便混用不同的锁对象。

// ┌────────────┐
// │  for cond  │◄────────────┐
// └────┬───────┘             │
//
//	   │True (条件不满足)     │
//	   ▼                     │
//	cond.Wait()              │
//	    │ (释放锁+阻塞)       │
//	    ▼                     │
//	[被唤醒 + 重新加锁]        │
//
//	    │                     │
//	    ▼                     │
//
// ┌────┴──────────────┐      │
// │ for 循环体结尾 → 回顶部 │───┘
// └───────────────────┘
//
//	   │
//	   ▼
//	条件满足 → 跳出循环
//	   │
//	   ▼
//	消费数据
// package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	n := 10
// 	m := 3
// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	cond := sync.NewCond(&mu)

// 	cur := 1
// 	turn := 0

// 	for i := 0; i < m; i++ {
// 		wg.Add(1)
// 		go func(id int) {
// 			defer wg.Done()
// 			for {
// 				mu.Lock()
// 				for id != turn && cur <= n {
// 					cond.Wait()
// 				}
// 				if cur > n {
// 					mu.Unlock()
// 					cond.Broadcast()
// 					return
// 				}
// 				fmt.Printf("goroutine%d: %d\n", id, cur)
// 				cur++
// 				turn = (turn + 1) % m
// 				cond.Broadcast()
// 				mu.Unlock()
// 			}
// 		}(i)
// 	}

// 	wg.Wait()
// }

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
	lock := sync.Mutex{}
	cond := sync.NewCond(&lock)
	cur := 0 //不能定义在for头部，因为for循环结束后cur会被销毁
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			for j := id; j <= m; j += n {
				lock.Lock()
				for j != cur {
					cond.Wait()
				}
				fmt.Println("当前goroutine为", id, "打印的数字为", cur)
				cur++
				lock.Unlock()
				cond.Broadcast()
			}
		}(i)
	}
	wg.Wait()
}

//-------------------------------------------------------------------
// 生产者-消费者问题的 Condition 版本
// 2 个消费者，1 个生产者，生产 5 个数据
// 注意：这个版本的代码更复杂一些，因为我们要确保消费者轮流消费。
// 这也是面试中常见的题目。

// package main

// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// )

// var (
// 	queue  = []int{}
// 	mu     = sync.Mutex{}
// 	cond   = sync.NewCond(&mu)
// 	wg     = sync.WaitGroup{}
// 	closed = false // 表示生产者是否已经结束
// 	turn   = 0     // 当前轮到哪个消费者
// )

// func consumer(id int) {
// 	defer wg.Done()
// 	for {
// 		mu.Lock()
// 		// 等待：当且仅当队列为空并且生产未结束时才等待
// 		for (len(queue) == 0 || turn != id) && !closed {
// 			cond.Wait()
// 		}

// 		// 如果队列为空并且生产已经结束，说明没有更多数据，安全退出
// 		if len(queue) == 0 && closed {
// 			mu.Unlock()
// 			return
// 		}

// 		// 否则从队列取数据处理
// 		x := queue[0]
// 		queue = queue[1:]
// 		mu.Unlock()
// 		turn = (turn + 1) % 2 // 切换到下一个消费者

// 		fmt.Println("Consumer", id, "got", x)
// 	}
// }

// func producer(n int) {
// 	for i := 1; i <= n; i++ {
// 		mu.Lock()
// 		queue = append(queue, i)
// 		cond.Signal() // 通知一个消费者
// 		fmt.Println("Produced", i)
// 		mu.Unlock()
// 		runtime.Gosched() // 让出时间片，让消费者有机会执行
// 	}
// 	// 生产结束，设置 closed 并唤醒所有等待者（否则可能有消费者永远阻塞）
// 	mu.Lock()
// 	closed = true
// 	cond.Broadcast()
// 	mu.Unlock()
// }

// func main() {
// 	for i := 0; i < 2; i++ {
// 		wg.Add(1)
// 		go consumer(i)
// 	}
// 	producer(5)
// 	wg.Wait() // 等待所有消费者退出
// }

//-------------------------------------------------------------------
//不使用runtime.Gosched()版本
// 这个版本中，生产者在每次生产后会调用 cond.Signal() 唤醒一个等待的消费者。
// 这样可以确保消费者有机会在生产者继续生产之前运行，从而避免生产者连续生产多个项目而消费者没有机会消费的情况。
// 这种方式更符合生产者-消费者模型的设计思路，避免了不必要的忙等待和资源浪费。

// package main

// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// )

// var (
// 	queue  = []int{}
// 	mu     = sync.Mutex{}
// 	condP  = sync.NewCond(&mu) // 控制生产者
// 	condC  = sync.NewCond(&mu) // 控制消费者
// 	wg     = sync.WaitGroup{}
// 	turn   = 0     // 当前轮到哪个消费者
// 	closed = false // 生产结束标志
// )

// func consumer(id int, totalConsumers int) {
// 	defer wg.Done()
// 	for {
// 		mu.Lock()
// 		for (len(queue) == 0 || turn != id) && !closed {
// 			fmt.Println("Consumer", id, "is waiting")
// 			condC.Wait() // 等待数据或等待轮到自己
// 		}

// 		if closed && len(queue) == 0 {
// 			mu.Unlock()
// 			return
// 		}

// 		x := queue[0]
// 		queue = queue[1:]
// 		fmt.Printf("Consumer %d got %d\n", id, x)

// 		// 切换到下一个消费者
// 		turn = (turn + 1) % totalConsumers

// 		// 通知生产者继续
// 		condP.Signal()
// 		mu.Unlock()
// 	}
// }

// func producer(n int, totalConsumers int) {
// 	for i := 1; i <= n; i++ {
// 		mu.Lock()
// 		for len(queue) > 0 { // 等待上一个被消费完
// 			condP.Wait()
// 		}

// 		queue = append(queue, i)
// 		fmt.Println("Produced", i) //实际上 fmt.Println("Produced", i) 确实被执行并输出了，
// 		// 可能的解释：
// 		// 输出缓冲问题：在某些环境下，标准输出可能被缓冲，导致输出没有立即显示
// 		// 竞争条件：在没有详细日志的原始版本中，可能存在微妙的竞争条件
// 		// 执行速度：生产者的执行可能太快，导致输出被后续的消费者输出覆盖
// 		condC.Signal() // 通知一个消费者
// 		mu.Unlock()
// 	}

// 	condC.Broadcast() // 唤醒所有消费者退出
// 	runtime.Gosched() // 让出时间片，让消费者有机会执行
// 	defer func() { closed = true }()
// }

// func main() {
// 	m := 2 // 两个消费者
// 	for i := 0; i < m; i++ {
// 		wg.Add(1)
// 		go consumer(i, m)
// 	}
// 	producer(5, m)
// 	wg.Wait()
// }

//-------------------------------------------------------------------
// 添加详细日志后，我们可以清楚地看到整个执行流程。
// goroutine调度不确定性
// goroutine调度器特性
// Go的运行时调度器负责在操作系统线程上调度goroutine的执行。调度决策是非确定性的，受多种因素影响：
// 系统负载
// 调度器状态
// 随机性因素
// CPU核心数量
// 这些因素共同作用，导致每次运行程序时，goroutine的执行顺序可能不同。
// 但是执行顺序不同的原因是goroutine调度的固有不确定性，这是正常且预期的行为。
// 程序在不同调度顺序下都能产生正确结果，这证明了同步逻辑的正确性。
// 这也是并发编程的一个重要目标：确保无论调度如何变化，程序都能正确运行。

// package main

// import (
// 	"fmt"
// 	"sync"
// )

// var (
// 	queue  = []int{}
// 	mu     = sync.Mutex{}
// 	condP  = sync.NewCond(&mu) // 控制生产者
// 	condC  = sync.NewCond(&mu) // 控制消费者
// 	wg     = sync.WaitGroup{}
// 	turn   = 0     // 当前轮到哪个消费者
// 	closed = false // 生产结束标志
// )

// func consumer(id int, totalConsumers int) {
// 	defer wg.Done()
// 	fmt.Printf("CONSUMER %d: Started\n", id)

// 	for {
// 		mu.Lock()
// 		fmt.Printf("CONSUMER %d: Lock acquired, queue: %v, turn: %d, closed: %v\n",
// 			id, queue, turn, closed)

// 		for (len(queue) == 0 || turn != id) && !closed {
// 			fmt.Printf("CONSUMER %d: Conditions not met - waiting\n", id)
// 			condC.Wait()
// 			fmt.Printf("CONSUMER %d: Woke up, queue: %v, turn: %d\n", id, queue, turn)
// 		}

// 		if closed && len(queue) == 0 {
// 			fmt.Printf("CONSUMER %d: Exiting\n", id)
// 			mu.Unlock()
// 			return
// 		}

// 		x := queue[0]
// 		queue = queue[1:]
// 		fmt.Printf("Consumer %d got %d\n", id, x)

// 		oldTurn := turn
// 		turn = (turn + 1) % totalConsumers
// 		fmt.Printf("CONSUMER %d: Turn changed from %d to %d\n", id, oldTurn, turn)

// 		fmt.Printf("CONSUMER %d: Signaling producer\n", id)
// 		condP.Signal()
// 		mu.Unlock()
// 		fmt.Printf("CONSUMER %d: Lock released\n", id)
// 	}
// }

// func producer(n int, totalConsumers int) {
// 	fmt.Println("PRODUCER: Starting producer")
// 	for i := 1; i <= n; i++ {
// 		fmt.Printf("PRODUCER: Attempting to produce %d\n", i)
// 		mu.Lock()
// 		fmt.Printf("PRODUCER: Lock acquired for %d, queue length: %d\n", i, len(queue))

// 		for len(queue) > 0 {
// 			fmt.Printf("PRODUCER: Waiting to produce %d, queue not empty\n", i)
// 			condP.Wait()
// 			fmt.Printf("PRODUCER: Woke up for %d, queue length: %d\n", i, len(queue))
// 		}

// 		queue = append(queue, i)
// 		fmt.Println("Produced", i)

// 		fmt.Printf("PRODUCER: Broadcasting to consumers for %d\n", i)
// 		condC.Broadcast()
// 		mu.Unlock()
// 		fmt.Printf("PRODUCER: Lock released for %d\n", i)
// 	}

// 	mu.Lock()
// 	closed = true
// 	fmt.Println("PRODUCER: Setting closed=true")
// 	condC.Broadcast()
// 	mu.Unlock()
// 	fmt.Println("PRODUCER: Finished")
// }

// func main() {
// 	m := 2 // 两个消费者
// 	for i := 0; i < m; i++ {
// 		wg.Add(1)
// 		go consumer(i, m)
// 	}
// 	producer(5, m)
// 	wg.Wait()
// }
