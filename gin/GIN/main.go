package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	//设置默认路由引擎
	r := gin.Default()

	//post主要用于提交数据
	r.POST("/post", func(c *gin.Context) {
		c.String(http.StatusOK, "This is a POST request")
	})
	//get主要用于获取数据
	r.GET("/get", func(c *gin.Context) {
		c.String(http.StatusOK, "This is a GET request")
	})
	//put主要用于更新数据
	r.PUT("/put", func(c *gin.Context) {
		c.String(200, "This is a PUT request")
	})
	//delete主要用于删除数据
	r.DELETE("/delete", func(c *gin.Context) {
		c.String(200, "This is a DELETE request")
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run(":8000")
}
