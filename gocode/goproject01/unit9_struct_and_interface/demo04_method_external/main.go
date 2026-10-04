package main

import (
	"fmt"

	"github.com/jhw66/unit9-module/model"
)

func main() {
	stu := model.Student{
		Name: "lll",
		Age:  99,
	}
	fmt.Println(stu)

	p := model.NewPerson("丽丽")
	p.SetAge(20)
	fmt.Println(stu.Name)
	//fmt.Println(p.name) 错误，因为p结构体中name开头无大写
	fmt.Println(p.GetAge())
	fmt.Println(p)
	fmt.Println(*p)
}
