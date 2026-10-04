package routers

import (
	"mygin/GIN05/controllers/nav"

	"github.com/gin-gonic/gin"
)

func NavGroup(r *gin.Engine) {
	navgroup := r.Group("/nav")
	{
		navgroup.POST("/read", nav.NavController{}.Read)
		navgroup.POST("update", nav.NavController{}.Update)
		navgroup.POST("/delete", nav.NavController{}.Delete)
		navgroup.POST("/add", nav.NavController{}.Add)
	}
}
