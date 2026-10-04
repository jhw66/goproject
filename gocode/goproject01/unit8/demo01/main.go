package main

import "fmt"

func main() {
	//使用make初始化
	var a map[int]string = make(map[int]string, 10) //map可以存放10个键值对
	//or a := make(map[int]string, 10) // map can store 10 key-value pairs
	//只声明map内存是没有分配空间的
	//必须通过make函数进行初始化，才会分配空间

	a[2009] = "hhh"
	a[2008] = "lll"
	a[2007] = "ppp"
	a[2009] = "hhh1" //会替代前面的
	fmt.Println(a)

	//使用字面量初始化
	b := map[int]string{
		2009: "hhh",
		2008: "lll",
	}
	b[2000] = "kkk"

	fmt.Print(b)

	//注意：
	//c:=map[int]int
	//无法直接使用，其实只是声明了一个 map 类型变量，它的值为 nil，并没有初始化底层的桶（bucket）。
	//
	//如果使用make初始化而不规定大小：
	//b := make(map[int]int)
	// 如果没有给容量 hint，Go 会默认初始化一个最小 bucket：
	// 默认 B = 0，也就是 2^0 = 1 个 bucket
	// 每个 bucket 可以存储 8 个 key/value 对
	// 内存大小大约 = hmap 头 + bucket 指针 + 1 个 bucket 的内存
	// 一个 bucket 的大小大约 = 8 * (key大小 + value大小 + metadata)
	// key 和 value 的大小取决于类型：
	// int32 → 4 字节
	// int64 → 8 字节
	// 所以最小 map 可以存储 8 个 key/value 对而不扩容。
	//
	//同理使用字面量初始化时：
	//b := map[int]int{
	//     1: 10,
	//     2: 20,
	// }
	// 编译器在 runtime 调用 makemap 时：
	// 根据字面量大小（len=2）计算 bucket 数量
	// 最小 bucket = 1 个（能容纳 8 个 key/value）
	// 所以对于小 map，字面量和 make(map[int]int) 内存分配几乎一样
}
