// package main

// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// )

// var counter int = 0

// func Count(lock *sync.Mutex, i int) {
// 	lock.Lock()
// 	counter++
// 	fmt.Println(counter, i)
// 	lock.Unlock()
// }
// func main() {
// 	lock := &sync.Mutex{}
// 	for i := 0; i < 10; i++ {
// 		go Count(lock, i)
// 	}
// 	for {
// 		lock.Lock()
// 		c := counter
// 		fmt.Println("hhh", c)
// 		lock.Unlock()
// 		runtime.Gosched()
// 		if c >= 10 {
// 			break
// 		}
// 	}
// }

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	shutdown int64
	wg       sync.WaitGroup
)

func main() {
	wg.Add(2)
	go doWork("A")
	go doWork("B")
	time.Sleep(1 * time.Second)
	fmt.Println("Shutdown Now")
	atomic.StoreInt64(&shutdown, 1)
	wg.Wait()
	fmt.Println(runtime.GOMAXPROCS(runtime.NumCPU()))
}
func doWork(name string) {
	defer wg.Done()
	for {
		fmt.Printf("Doing %s Work\n", name)
		time.Sleep(250 * time.Millisecond)
		if atomic.LoadInt64(&shutdown) == 1 {
			fmt.Printf("Shutting %s Down\n", name)
			break
		}
	}
}

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"sync/atomic"
// )

// // 版本1：使用原子操作（正确）
// func atomicVersion() {
// 	var counter int64
// 	var wg sync.WaitGroup

// 	for i := 0; i < 1000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			atomic.AddInt64(&counter, 1)
// 			wg.Done()
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Printf("Atomic result: %d\n", counter)
// }

// // 版本2：不使用原子操作（有竞态条件）
// func nonAtomicVersion() {
// 	var counter int64
// 	var wg sync.WaitGroup

// 	for i := 0; i < 10000; i++ {
// 		wg.Add(1)
// 		go func() {
// 			counter++ // 非原子操作
// 			wg.Done()
// 		}()
// 	}

// 	wg.Wait()
// 	fmt.Printf("Non-atomic result: %d\n", counter)
// }

// func main() {
// 	// 多次运行，观察结果
// 	for i := 0; i < 5; i++ {
// 		atomicVersion()    // 总是1000
// 		nonAtomicVersion() // 可能小于1000
// 		fmt.Println()
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// 	"sync/atomic"
// )

// var (
// 	counter int64
// 	wg      sync.WaitGroup
// )

// func main() {
// 	wg.Add(2)
// 	go incCounter(1)
// 	go incCounter(2)
// 	wg.Wait() //等待goroutine结束
// 	fmt.Println(counter)
// 	fmt.Println(runtime.NumCPU())

// }

// func incCounter(id int) {
// 	defer wg.Done()
// 	for count := 0; count < 20000; count++ {
// 		atomic.AddInt64(&counter, 1) //安全的对counter加1
// 		runtime.Gosched()
// 	}
// }
