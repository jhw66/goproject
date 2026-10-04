// demo02：WithCancel / defer cancel / 取消传播。
// 关键点：
//  1. 每个 WithXxx 返回的 cancel 必须在所有路径调用（通常 defer），否则子节点与定时器会泄漏。
//  2. 取消会沿派生链向下传播：父取消后，所有由它派生的子 ctx 也会立刻取消。
//  3. 主动取消时 ctx.Err() 为 context.Canceled（errors.Is 判断）。
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			fmt.Printf("[%s] 退出: err=%v, isCanceled=%v\n",
				name, err, errors.Is(err, context.Canceled))
			return
		case <-time.After(120 * time.Millisecond):
			fmt.Printf("[%s] 心跳\n", name)
		}
	}
}

func main() {
	parent, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()

	child1, cancelChild1 := context.WithCancel(parent)
	defer cancelChild1()

	child2, cancelChild2 := context.WithCancel(parent)
	defer cancelChild2()

	child3, cancelChild3 := context.WithCancel(parent)
	defer cancelChild3()

	child4, cancelChild4 := context.WithCancel(child3)
	defer cancelChild4()

	go worker(parent, "parent")
	go worker(child1, "child1")
	go worker(child2, "child2")
	go worker(child3, "child3")
	go worker(child4, "child4")

	time.Sleep(1 * time.Second)
	fmt.Println(">>> 取消 child3，观察 child4 是否一并退出")
	cancelChild3()

	time.Sleep(1 * time.Second)
	fmt.Println(">>> 取消 child2")
	cancelChild2()

	time.Sleep(3 * time.Second)
	fmt.Println(">>> 取消 parent，观察 child1 是否一并退出")
	cancelParent()

	time.Sleep(150 * time.Millisecond)
}
