//只用 sync.Mutex 实现严格有序打印的版本

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	n := 50
	m := 4

	var mu sync.Mutex
	cur := 1 // 当前要打印的数字
	var wg sync.WaitGroup

	for i := 0; i < m; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				mu.Lock()
				if cur > n {
					mu.Unlock()
					return
				}
				if (cur-1)%m == id {
					fmt.Printf("Goroutine %d received: %d\n", id, cur)
					cur++
				}
				mu.Unlock()

				// 稍作休息，防止其他 goroutine 抢锁过快造成忙等
				time.Sleep(time.Microsecond * 10)
				//给其他 goroutine 一点执行机会，防止某个 goroutine一直占着 CPU 反复抢锁，从而导致忙等（busy waiting）现象。
			}
		}(i)
	}

	wg.Wait()
}
