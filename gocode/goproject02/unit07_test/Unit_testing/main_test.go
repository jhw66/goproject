package main

import (
	"testing"
	"time"
)

func TestFibonacci(t *testing.T) {
	fsMap := map[int]int{}

	fsMap[-1] = 0
	fsMap[0] = 0
	fsMap[1] = 1
	fsMap[2] = 1
	fsMap[3] = 2
	fsMap[4] = 3
	fsMap[5] = 5
	fsMap[6] = 8
	fsMap[7] = 13
	fsMap[8] = 21
	fsMap[9] = 34

	time.Sleep(5 * time.Second)

	for k, v := range fsMap {
		fib := Fibonacci(k)
		if v == fib {
			t.Logf("结果正确,Fibonacci(%d)=%d", k, fib)
		} else {
			t.Errorf("结果错误,Fibonacci(%d)=%d,expected %d", k, fib, v)
		}
	}
}

func TestFibonacci2(t *testing.T) {
	fsMap := map[int]int{}

	fsMap[-1] = 0
	fsMap[0] = 0
	fsMap[1] = 1
	fsMap[2] = 1
	fsMap[3] = 2
	fsMap[4] = 3
	fsMap[5] = 5
	fsMap[6] = 7
	fsMap[7] = 13
	fsMap[8] = 21
	fsMap[9] = 34

	for k, v := range fsMap {
		fib := Fibonacci(k)
		if v == fib {
			t.Logf("结果正确,Fibonacci(%d)=%d", k, fib)
		} else {
			t.Errorf("结果错误,Fibonacci(%d)=%d,expected %d", k, fib, v)
		}
	}
}

func BenchmarkFibonacci(b *testing.B) {
	b.ReportAllocs() //报告内存分配
	b.ResetTimer()   //重置计时器
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}
