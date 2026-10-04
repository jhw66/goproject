// package main

// import (
// 	"fmt"
// 	"runtime"
// 	"sync"
// 	"time"
// )

// type value struct {
// 	memAccess sync.Mutex
// 	value     int
// }

// func main() {
// 	runtime.GOMAXPROCS(3)
// 	var wg sync.WaitGroup
// 	sum := func(v1, v2 *value) {
// 		defer wg.Done()
// 		v1.memAccess.Lock()
// 		time.Sleep(2 * time.Second)
// 		v2.memAccess.Lock()
// 		fmt.Printf("sum = %d\n", v1.value+v2.value)
// 		v2.memAccess.Unlock()
// 		v1.memAccess.Unlock()
// 	}
// 	product := func(v1, v2 *value) {
// 		defer wg.Done()
// 		v2.memAccess.Lock()
// 		time.Sleep(2 * time.Second)
// 		v1.memAccess.Lock()
// 		fmt.Printf("product = %d\n", v1.value*v2.value)
// 		v1.memAccess.Unlock()
// 		v2.memAccess.Unlock()
// 	}
// 	var v1, v2 value
// 	v1.value = 1
// 	v2.value = 1
// 	wg.Add(2)
// 	go sum(&v1, &v2)
// 	go product(&v1, &v2)
// 	wg.Wait()
// }

// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"runtime"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// func main() {
// 	runtime.GOMAXPROCS(3)
// 	cv := sync.NewCond(&sync.Mutex{})
// 	go func() {
// 		for range time.Tick(1 * time.Second) { // 通过tick控制两个人的步调
// 			cv.Broadcast()
// 		}
// 	}()
// 	takeStep := func() {
// 		cv.L.Lock()
// 		cv.Wait()
// 		cv.L.Unlock()
// 	}
// 	tryDir := func(dirName string, dir *int32, out *bytes.Buffer) bool {
// 		fmt.Fprintf(out, " %+v", dirName)
// 		atomic.AddInt32(dir, 1)
// 		takeStep()                      //走上一步
// 		if atomic.LoadInt32(dir) == 1 { //走成功就返回
// 			fmt.Fprint(out, ". Success!")
// 			return true
// 		}
// 		takeStep() // 没走成功，再走回来
// 		atomic.AddInt32(dir, -1)
// 		return false
// 	}
// 	var left, right int32
// 	tryLeft := func(out *bytes.Buffer) bool {
// 		return tryDir("向左走", &left, out)
// 	}
// 	tryRight := func(out *bytes.Buffer) bool {
// 		return tryDir("向右走", &right, out)
// 	}
// 	walk := func(walking *sync.WaitGroup, name string) {
// 		var out bytes.Buffer
// 		defer walking.Done()
// 		defer func() { fmt.Println(out.String()) }()
// 		fmt.Fprintf(&out, "%v is trying to scoot:", name)
// 		for i := 0; i < 5; i++ {
// 			if tryLeft(&out) || tryRight(&out) {
// 				return
// 			}
// 		}
// 		fmt.Fprintf(&out, "\n%v is tried!", name)
// 	}
// 	var trail sync.WaitGroup
// 	trail.Add(2)
// 	go walk(&trail, "男人") // 男人在路上走
// 	go walk(&trail, "女人") // 女人在路上走
// 	trail.Wait()
// }

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(3)
	var wg sync.WaitGroup
	const runtime = 1 * time.Second
	var sharedLock sync.Mutex
	greedyWorker := func() {
		defer wg.Done()
		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}
		fmt.Printf("Greedy worker was able to execute %v work loops\n", count)
	}
	politeWorker := func() {
		defer wg.Done()
		var count int
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops\n", count)
	}
	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
}
