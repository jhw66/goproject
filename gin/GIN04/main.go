package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Age      int    `json:"age" form:"age"`
}

type XmlData struct {
	Title   string `json:"title" xml:"title"`
	Content string `json:"content" xml:"content"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*")

	//Get 请求传值
	r.GET("/", func(c *gin.Context) {
		c.String(200, "首页")

		username := c.Query("username")
		age := c.DefaultQuery("age", "18")

		c.JSON(200, gin.H{
			"username": username,
			"age":      age,
		})
	})

	//获取 Get 请求传值，绑定到结构体
	r.GET("/getUserInfo", func(c *gin.Context) {
		userInfo := UserInfo{}
		if err := c.ShouldBind(&userInfo); err == nil {
			fmt.Printf("%#v\n", userInfo)
			c.JSON(200, userInfo)
			c.String(200, "hhh")
		} else {
			c.JSON(200, gin.H{
				"error": err.Error(),
			})
		}
	})

	//Post 表单提交
	r.GET("/form", func(c *gin.Context) {
		c.HTML(200, "form.html", nil)
	})
	//获取 form 表单提交的数据
	r.POST("/doAddUser", func(c *gin.Context) {
		// username := c.PostForm("username")
		// password := c.PostForm("password")
		// age := c.DefaultPostForm("age", "20")
		// c.JSON(200, gin.H{
		// 	"username": username,
		// 	"password": password,
		// 	"age":      age,
		// })
		//绑定到结构体
		userInfo := UserInfo{}
		if err := c.ShouldBind(&userInfo); err == nil {
			c.JSON(200, userInfo)
		} else {
			c.JSON(200, gin.H{
				"error": err.Error(),
			})
		}
	})

	//获取xml数据
	r.POST("/xml", func(c *gin.Context) {
		// 		<?xml version="1.0" encoding="UTF-8"?>
		// <article>
		//     <title type="string">wjh666</title>
		//     <content type="string">666</content>
		// </article>
		xmlData := &XmlData{}
		xmlSliceData, _ := c.GetRawData() //获取请求体数据
		if err := xml.Unmarshal(xmlSliceData, xmlData); err == nil {
			c.JSON(200, xmlData)
			fmt.Println(string(xmlSliceData))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

	})

	//动态路由
	r.GET("/user/:name/:age", func(c *gin.Context) {
		name := c.Param("name")
		age := c.Param("age")
		c.JSON(200, gin.H{
			"name": name,
			"age":  age,
		})
	})

	r.Run()
}
