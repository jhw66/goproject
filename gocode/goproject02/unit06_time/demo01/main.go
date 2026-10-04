package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

	now := time.Now()
	fmt.Println(now)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := now.In(loc)
	fmt.Println(shanghaiTime)

}
