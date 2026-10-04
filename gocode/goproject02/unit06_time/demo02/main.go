package main

import (
	"fmt"
	"time"
)

// -------------------------------------------------------
// 1. time.Sleep — 最简单的延迟
// -------------------------------------------------------
func demoSleep() {
	fmt.Println("\n===== 1. time.Sleep =====")
	fmt.Println("开始 Sleep，等待 1 秒...")
	start := time.Now()
	time.Sleep(1 * time.Second)
	fmt.Printf("Sleep 结束，实际耗时: %v\n", time.Since(start).Round(time.Millisecond))
}

// -------------------------------------------------------
// 2. time.After — 一次性延迟（select 超时控制）
// -------------------------------------------------------
func demoAfter() {
	fmt.Println("\n===== 2. time.After =====")

	ch := make(chan string)

	// 模拟一个耗时任务：500ms 后完成
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "任务完成"
	}()

	select {
	case result := <-ch:
		fmt.Println("成功:", result)
	case <-time.After(2 * time.Second):
		fmt.Println("超时！")
	}

	// 模拟超时场景：任务耗时 3s，超时阈值 1s
	go func() {
		time.Sleep(3 * time.Second)
		ch <- "慢任务完成"
	}()

	select {
	case result := <-ch:
		fmt.Println("成功:", result)
	case <-time.After(1 * time.Second):
		fmt.Println("超时！慢任务未能在 1 秒内完成")
	}
}

// -------------------------------------------------------
// 3. time.NewTimer — 一次性定时器（完整版，支持 Stop/Reset）
// -------------------------------------------------------
func demoNewTimer() {
	fmt.Println("\n===== 3. time.NewTimer =====")

	// 3-1 正常触发
	timer1 := time.NewTimer(5 * time.Second)
	<-timer1.C
	fmt.Println("timer1 触发（1s）")

	// 3-2 提前 Stop
	timer2 := time.NewTimer(5 * time.Second)
	stopped := timer2.Stop()
	if stopped {
		fmt.Println("timer2 被提前取消，未触发")
	}

	// 3-3 Reset 重新计时
	timer3 := time.NewTimer(2 * time.Second)
	// 先 Stop 并排空 channel
	if !timer3.Stop() {
		select {
		case <-timer3.C:
		default:
		}
	}
	timer3.Reset(5 * time.Second)
	<-timer3.C
	fmt.Println("timer3 Reset 后重新触发（5s）")

	// 3-4 在 goroutine 中取消定时器
	timer4 := time.NewTimer(3 * time.Second)
	quit := make(chan struct{})

	go func() {
		time.Sleep(200 * time.Millisecond)
		if !timer4.Stop() {
			<-timer4.C
		}
		fmt.Println("timer4 在 goroutine 中被取消")
		close(quit)
	}()

	select {
	case <-timer4.C:
		fmt.Println("timer4 触发（不应到达这里）")
	case <-quit:
		fmt.Println("timer4 已通过 quit 信号确认取消")
	}
}

// -------------------------------------------------------
// 4. time.NewTicker — 周期性计时器
// -------------------------------------------------------
func demoNewTicker() {
	fmt.Println("\n===== 4. time.NewTicker =====")

	// 4-1 基本用法：固定次数后停止
	ticker1 := time.NewTicker(500 * time.Millisecond)
	count := 0
	for t := range ticker1.C {
		count++
		fmt.Printf("  ticker1 第 %d 次: %s\n", count, t.Format("15:04:05.000"))
		if count >= 4 {
			ticker1.Reset(200 * time.Millisecond)
		}
		if count >= 8 {
			ticker1.Stop()
			break
		}
	}

	// 4-2 结合 select 与退出信号
	fmt.Println("--- ticker2：结合 select 退出 ---")
	ticker2 := time.NewTicker(200 * time.Millisecond)
	done := time.After(1 * time.Second) // 1 秒后退出
	for {
		select {
		case t := <-ticker2.C:
			fmt.Printf("  ticker2 tick: %s\n", t.Format("15:04:05.000"))
		case <-done:
			ticker2.Stop()
			fmt.Println("  ticker2 退出")
			return
		}
	}
}

// -------------------------------------------------------
// 5. time.AfterFunc — 异步延迟执行（新 goroutine 中运行）
// -------------------------------------------------------
func demoAfterFunc() {
	fmt.Println("\n===== 5. time.AfterFunc =====")

	wait := make(chan struct{})

	// 5-1 正常触发
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("AfterFunc 触发（5s 后在新 goroutine 中执行）")
		close(wait)
	})
	<-wait // 等待回调完成

	// 5-2 提前取消
	timer := time.AfterFunc(5*time.Second, func() {
		fmt.Println("AfterFunc 触发（不应到达这里）")
	})
	time.Sleep(2 * time.Second)
	stopped := timer.Stop()
	if stopped {
		fmt.Println("AfterFunc 定时器被提前取消")
	}
}

// -------------------------------------------------------
// main
// -------------------------------------------------------
func main() {
	demoSleep()
	demoAfter()
	demoNewTimer()
	demoNewTicker()
	demoAfterFunc()

	fmt.Println("\n===== 全部示例执行完毕 =====")
}
