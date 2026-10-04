package main

//方法是作用在指定的数据类型上，和指定的数据类型绑定，因此自定义类型都可以有方法，不止struct
import "fmt"

type A struct {
	Num  int
	Name string
}

// 给A结构体绑定方法test
// 假如方法名首字母大写，就可以在本包和其他包访问
func (a A) test() {
	fmt.Println("调用A方法test", a.Name, a.Num)
	a.Name = "HHH"
	fmt.Println("改变Name", a.Name)
}
func (a *A) test2() {
	a.Num = 7
}
func (s A) String() string {
	str := fmt.Sprintf("Name=%v,Age=%v,调用了String方法", s.Name, s.Num)
	return str
}
func (s *A) string() string {
	str := fmt.Sprintf("Name=%v,Age=%v,调用了String方法", s.Name, s.Num)
	return str
}

// 自定义类型绑定方法
type integer int

func (i integer) print() {
	fmt.Println("i=", i)
}

func main() {
	var p A = A{88, "HHH"}
	fmt.Println(p)
	p.test()
	p.Name = "hhh"
	p.test() //跟普通函数一样是值传递
	fmt.Println(p)
	p.test2()
	fmt.Println(p)

	var i integer = 20
	i.print()

	stu := A{
		Name: "qqq",
		Num:  30,
	}
	//传入地址，假如绑定了String方法,在fmt.Println就会自动调用
	fmt.Println(stu)
	fmt.Println(&stu)
	fmt.Println(stu.string())
	fmt.Println((&stu).string())
}

//  String() 实现的是 fmt.Stringer 接口
//  go
// // 这是 Go 内置的接口定义
// type Stringer interface {
//     String() string
// }
// 自动接口实现机制：
// 当你定义了一个 String() string 方法时，Go 编译器会自动识别你的类型实现了 Stringer 接口
// 这叫做隐式接口实现 - Go 不需要显式声明实现了哪个接口
// fmt 包中的函数会检查传入的值是否实现了 Stringer 接口，如果是，就调用 String() 方法
// go
// package fmt
// // 简化版的实现逻辑
// func Println(a ...interface{}) {
//     for _, arg := range a {
//         if s, ok := arg.(fmt.Stringer); ok {
//             // 如果实现了 Stringer 接口，调用 String() 方法
//"这里是不是就是看arg接口是否可以转成Stringer类型接口，如果可以转就说明arg接口实现了String方法就可以直接调用"
//             writeString(s.String())
//         } else {
//             // 否则使用默认的格式化
//             defaultFormat(arg)
//         }
//     }
// }

// // 1. 值接收者 - 非常灵活
// func (s A) String() string {
//     return fmt.Sprintf("值接收者: %s-%d", s.Name, s.Num)
// }

// func main() {
//     // 所有情况都可用：
//     stu1 := A{Name: "张三", Num: 20}
//     stu2 := &A{Name: "李四", Num: 25}

//     fmt.Println(stu1)    // ✅ 值调用值接收者方法
//     fmt.Println(stu2)    // ✅ 指针调用值接收者方法（自动解引用）
//     fmt.Println(A{Name: "王五", Num: 30}) // ✅ 临时值也能调用
// }

// // 2. 指针接收者 - 有局限性
// func (s *A) String() string {
//     return fmt.Sprintf("指针接收者: %s-%d", s.Name, s.Num)
// }

// func main() {
//     stu1 := A{Name: "张三", Num: 20}
//     stu2 := &A{Name: "李四", Num: 25}

//     fmt.Println(&stu1)   // ✅ 可以调用（传入指针）
//     fmt.Println(stu1)    // ⚠️ 可能不行！如果stu1不可寻址会出错
//     fmt.Println(stu2)    // ✅ 可以调用（传入指针）
//     fmt.Println(A{Name: "王五", Num: 30}) // ❌ 编译错误！临时值不可取地址
// }
