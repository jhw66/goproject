package admin

import (
	models "mygin/GIN05/models"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

var db = models.DB

func (con BaseController) success(c *gin.Context) {
	c.String(200, "成功")
}

func (con BaseController) fail(c *gin.Context) {
	c.String(200, "失败")
}
