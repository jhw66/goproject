package main

import (
	router "mygin/GIN05/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	//gin.Default已经有中间件：engine.Use(Logger(), Recovery())
	//Logger中间件把日志写入gin.DefaultWriter,即使配置了GIN_MODE=relese
	//Recovery中间件会recover任何panic，如果有panic的话，会写入500响应码
	//如果不想使用默认中间件，可以使用gin.New()来创建路由引擎
	r.Static("/static", "./static")

	router.DefaultGroup(r)
	router.ApiGroup(r)
	router.AdminGroup(r)
	router.CrudGroup(r)
	router.NavGroup(r)
	router.PeopleGroup(r)

	r.Run(":80")

}
