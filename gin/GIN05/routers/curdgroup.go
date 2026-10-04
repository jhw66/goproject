package routers

import (
	"mygin/GIN05/controllers/curd"

	"github.com/gin-gonic/gin"
)

func CrudGroup(r *gin.Engine) {
	curdGroup := r.Group("/curd")
	{
		curdGroup.POST("/add", curd.CrudController{}.Add)
		curdGroup.POST("/read", curd.CrudController{}.Read)
		curdGroup.POST("/delete", curd.CrudController{}.Delete)
		curdGroup.POST("/update", curd.CrudController{}.Update)
	}
}
