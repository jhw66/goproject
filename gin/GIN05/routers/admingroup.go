package routers

import (
	"mygin/GIN05/controllers/admin"

	"github.com/gin-gonic/gin"
)

func AdminGroup(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	//配置全局中间件
	//adminGroup.Use(middlewares.Middleini())
	{
		adminGroup.GET("/", func(c *gin.Context) {
			c.String(200, "后台管理页面")
		})
		//注意函数不能加(),因为加了()表示函数调用
		//当访问路由时函数才会调用
		adminGroup.GET("/user", admin.UserController{}.Index)
		adminGroup.GET("/user/add", admin.UserController{}.Addfiles)
		adminGroup.POST("/user/doUpload", admin.UserController{}.DoUpload)
		adminGroup.GET("/user/edit", admin.UserController{}.Edit)
		adminGroup.GET("/order", func(c *gin.Context) {
			c.String(200, "后台新闻列表页面")
		})
		adminGroup.GET("/order/add", func(c *gin.Context) {
			c.String(200, "后台订单管理-add")
		})
		adminGroup.GET("/order/edit", func(c *gin.Context) {
			c.String(200, "后台订单管理-edit")
		})
		//adminGroup.GET("/student", admin.Student{}.Index)
		//adminGroup.GET("/lesson", admin.Lesson{}.Index)
		adminGroup.GET("/bank", admin.BankController{}.Index)
		adminGroup.POST("/bankcrud", admin.BankController{}.CRUD)
	}

}
