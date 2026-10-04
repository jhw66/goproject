package main

import "fmt"

// 结构体定义
type Teacher struct {
	Name   string
	Age    int
	School string
}

func main() {
	//创建对象
	var ma Teacher
	fmt.Println(ma)
	ma.Name = "ma"
	ma.Age = 22
	ma.School = "qinghua"
	fmt.Println(ma)

	var ha Teacher = Teacher{"ha", 23, "beida"}
	fmt.Println(ha)

	var t *Teacher = new(Teacher)
	(*t).Name = "sa"
	(*t).Age = 45
	t.School = "zhejiangu" //底层转化为(*t).School
	fmt.Println(*t)

	var r *Teacher = &Teacher{"rr", 11, "fudau"}
	fmt.Println(r, " ", *r)

	var l Teacher = Teacher{
		Name: "ll",
		Age:  33,
	}
	fmt.Println(l)
}

// type Stu struct {
//     Name string
//     Age  int
// }
// 1. var s Stu
// 定义一个 Stu 类型的变量，值是零值（字段均为默认零值）。
// var s Stu
// fmt.Println(s.Name, s.Age) // "" 0
// 特点：
// 分配在栈或堆由编译器决定
// 不是指针
// 初始内容：全部字段是零值

// 2. s := Stu
// 这个写法根本不能用，会编译错误：
// non-name Stu on left side of :=

// 3. s := Stu{}
// 和 var s Stu 效果等价，都是创建一个零值的结构体。
// 特点：
// 也是值类型
// 字段均为零值
// 语义清晰，常用于初始化

// 4. s := new(Stu)
// 创建一个结构体的 指针，字段为零值。
// s := new(Stu)
// fmt.Printf("%T\n", s) // *Stu
// 特点：
// 返回的是 *Stu
// 和 &Stu{} 等价：
// s := &Stu{}

// 5. make(Stu)
// 错误用法 —— make 只能用于 slice、map、chan。
// 结构体不是 make 的对象：
// cannot make Stu
// 正确使用：
// make([]int, 0)
// make(map[string]int)
// make(chan int)
