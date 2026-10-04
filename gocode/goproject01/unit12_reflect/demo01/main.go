package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func testReflect(i interface{}) {
	//调用TypeOf函数，返回reflect，reflect.Type类型数据
	reType := reflect.TypeOf(i)

	//调用ValueOf函数，返回reflect，reflect.Value类型数据
	reValue := reflect.ValueOf(i)

	fmt.Println(reType, reValue, reType.Kind())
	fmt.Printf("reType的数据类型为:%T,reValue的数据类型为:%T\n", reType, reValue)

	//如果想要获取reValue的数值，要调用Int()方法：返回v持有的有符号整数
	num2 := 80 + reValue.Int()
	fmt.Println(num2)

	//转回去
	i2 := reValue.Interface()
	//利用断言
	n := i2.(int)
	fmt.Println(n)
}
func main() {
	var num int = 100
	//将基本数据类型转到一个空接口才可以进行反射
	testReflect(num)

	stu := Student{
		Name: "hh",
		Age:  18,
	}
	//获取结构体的值,Elem()方法获取指针指向的值
	value := reflect.ValueOf(&stu).Elem()
	fmt.Println(value)

	//修改结构体的值
	name := value.FieldByName("Name")
	name.SetString("hhh")
	fmt.Println(value)

	//获取结构体的类型
	Type := reflect.TypeOf(&stu).Elem()
	fmt.Println(Type)
	fmt.Println(Type.Kind())
	fmt.Println(Type.Name())
	fmt.Println(Type.String())
	fmt.Println(value.Kind())

	//获取结构体的字段
	for i := 0; i < Type.NumField(); i++ {
		field := Type.Field(i)
		v := value.Field(i)
		fmt.Printf("第%d个字段的名字是:%s,字段类型是:%s,字段的值是:%v\n", i, field.Name, field.Type, v.Interface())
	}
}
