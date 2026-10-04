// Go 语言中**官方定义的引用类型（reference types）**有 4 种：

// 类型	示例	底层结构中包含的“引用（指针）”
// slice	[]int	指向底层数组的指针
// map	map[string]int	指向底层哈希表结构的指针
// channel	chan int	指向底层 channel 结构的指针
// function	func()	指向函数及其闭包环境的指针

// 有时也把 interface 也算作“间接引用”类型（因为它内部也是指针结构），但在官方语义里，interface 是一种独立的类型系统机制，
// 不是正式归类在 reference types 里。

// 为什么叫“引用类型”
// Go 中所有赋值、传参都是值传递（copy）。
// 但对于这些“引用类型”，这个“值”本身内部包含指针，所以 copy 的其实是指针结构。
// 因此，多个变量可以共享底层数据。
package main

import (
	"fmt"
)

// // Slice
// // 内部结构（Go runtime 定义在 runtime/slice.go）：
// type slice struct {
// 	array unsafe.Pointer // 底层数组的指针
// 	len   int
// 	cap   int
// }

// // 内部结构在 runtime/map.go，是一个指针指向哈希表对象的包装：
// type mapType struct {
// 	hmap *hmap // 指向底层哈希表结构
// }
// type hmap struct {
// 	count int
// 	B     uint8
// 	//...省略若干字段...
// }

// // 内部实现定义在 runtime/chan.go，大致结构是：
// type hchan struct {
// 	qcount uint
// 	buf    unsafe.Pointer // 环形缓冲区
// 	sendx  uint
// 	recvx  uint
// 	//...
// }

// // 函数值在 Go 中也是一个结构体（runtime 里定义为 funcval）：
// type funcval struct {
// 	fn  uintptr        // 函数入口地址
// 	ptr unsafe.Pointer // 捕获的闭包环境（如果有）
// }

func main() {
	a := []int{1, 2, 3}
	b := a         // 这里 a 和 b 拷贝了 slice header，但指向同一个底层数组。
	b[0] = 99      // 修改底层数组
	fmt.Println(a) // [99 2 3]

	m1 := map[string]int{"x": 1}
	m2 := m1 // 这里 m1、m2 都指向同一个底层哈希表结构。
	m2["x"] = 99
	fmt.Println(m1["x"]) // 99

	//ch 实际上是一个指向 hchan 的指针
	ch := make(chan int, 1)
	ch2 := ch //它们操作的是同一个底层通道。
	ch2 <- 1
	fmt.Println(<-ch) // 1

	x := 10
	f := func() { fmt.Println(x) }
	g := f //f 和 g 拷贝的函数值结构体内部，包含了指向闭包环境的指针
	g()    // 打印 10

}

// Go 中的“引用类型”指的是：
// 其值内部包含指向底层数据结构的指针，因此在赋值或传参时，多个变量可以共享同一份底层数据。
// 分类	    值类型（value types）	     引用类型（reference types）
// 例子  	int、float、array、struct	 slice、map、chan、func
// 行为  	赋值拷贝完整内容	          赋值只拷贝指针结构
// 共享数据	否	                         是
