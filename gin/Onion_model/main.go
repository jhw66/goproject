package main

import "fmt"

type Context struct {
	handlers []func(c *Context)
	position int64
}

func (this *Context) Use(f func(c *Context)) {
	this.handlers = append(this.handlers, f)
}
func (this *Context) Get(path string, f func(c *Context)) {
	this.handlers = append(this.handlers, f)
}
func (this *Context) Next() {
	this.position++
	this.handlers[this.position](this)
}
func (this *Context) Run() {
	this.position = 0
	this.handlers[this.position](this)
}

func main() {
	c := &Context{}
	c.Use(middleware1())
	c.Use(middleware2())
	c.Get("/", func(c *Context) {
		fmt.Println("Get handlerfunc  run")
	})
	c.Run()
}

func middleware1() func(c *Context) {
	return func(c *Context) {
		fmt.Println("middleware1 begin")
		c.Next()
		fmt.Println("middleware1 end")
	}
}

func middleware2() func(c *Context) {
	return func(c *Context) {
		fmt.Println("middleware2 begin")
		c.Next()
		fmt.Println("middleware2 end")
	}
}
