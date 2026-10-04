package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Middleini() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, time.Now().Format("2006-01-02 15:04:05")+c.Request.URL.String())
		c.Set("username", "张三")
	}
}
