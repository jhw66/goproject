package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Article struct {
	//字段首字母必须大写，否则无法导出，JSON无法解析
	Status  string `json:"status"`
	Message string `json:"message"`
	Numbers int    `json:"numbers"`
}

func main() {
	r := gin.Default()
	//加载模板文件
	r.LoadHTMLGlob("template/*")
	r.GET("/", func(c *gin.Context) {
		c.String(200, "首页%v", "wjh")
	})
	r.GET("/json01", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"message": "hello wjh",
			"numbers": 123,
		})
	})
	r.GET("/json02", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "hello wjh",
			"numbers": 123,
		})
	})
	r.GET("/json03", func(c *gin.Context) {
		a := &Article{
			Status:  "ok",
			Message: "hello wjh",
			Numbers: 123,
		}
		c.JSON(200, a)
	})
	//jsonp主要用于跨域请求
	//http://localhost:8080/jsonp?callback=xxx
	//xxx是函数名
	//返回结果是xxx({...})
	r.GET("/jsonp", func(c *gin.Context) {
		data := gin.H{
			"status":  "ok",
			"message": "hello wjh",
			"numbers": 123,
		}
		c.JSONP(200, data)
	})

	r.GET("/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "hello wjh",
			"numbers": 123,
		})
	})

	r.GET("/html01", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods.html", gin.H{
			"price": []int{100, 200, 300, 400, 500},
		})
	})

	r.GET("/html02", func(c *gin.Context) {
		c.HTML(200, "news.html", gin.H{
			"name": "wjh",
			"context": &Article{
				Status:  "ok",
				Message: "hello wjh",
				Numbers: 123,
			},
		})
	})
	r.Run()
}
