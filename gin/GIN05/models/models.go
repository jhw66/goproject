package models

import "time"

//如果我们一个功能想在多个控制器，或者多个模板里面复用的话，
//那么外面就可以把公共的功能单独抽取出来作为一个模块(Model)

func GetDate() string {
	tem := "20060102"
	return time.Now().Format(tem)
}

func GetUnix() int64 {
	return time.Now().Unix()
}
