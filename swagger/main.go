package main

import (
	_ "swagger/docs" // 重要：导入自动生成的docs包

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 用户管理API
// @version 1.0
// @description 这是一个用户管理的API文档
// @termsOfService http://example.com/terms/
// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	r := gin.Default()

	// 配置Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 您的API路由配置
	api := r.Group("/api/v1")
	{
		api.GET("users/:id", GetUserByID)
		api.POST("users", CreateUser)
	}

	r.Run(":8080")
}

// User 用户模型
// @Description 用户基本信息
type User struct {
	ID       int    `json:"id" example:"1"`                       // 用户ID
	Username string `json:"username" example:"zhangsan"`          // 用户名
	Email    string `json:"email" example:"zhangsan@example.com"` // 用户邮箱
	Age      int    `json:"age" example:"25"`                     // 用户年龄
}

// ErrorResponse 错误响应
// @Description 接口错误时的返回信息
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"请求参数错误"`
}

// GetUserByID
// @Summary 获取指定用户信息
// @Description 根据用户ID获取用户信息
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param id path int true "用户ID" minimum(1)
// @Param x-token header string true "认证Token"
// @Success 200 {object} User "成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 404 {object} ErrorResponse "未找到用户"
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	// 处理逻辑
	user := User{
		ID:       1,
		Username: "wjh",
		Email:    "wjh@qq.com",
		Age:      20,
	}
	c.JSON(200, user)
}

// CreateUser
// @Summary 创建用户
// @Description 创建新的用户
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param user body User true "用户信息"
// @Success 201 {object} User "创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部服务器错误"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 创建用户的逻辑...
	c.JSON(201, user)
}
