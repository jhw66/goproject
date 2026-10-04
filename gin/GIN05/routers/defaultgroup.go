package routers

import (
	"mygin/GIN05/controllers/itying"

	"github.com/gin-gonic/gin"
)

func DefaultGroup(r *gin.Engine) {
	r.LoadHTMLGlob("template/*") //此处template前不能使用..为什么
	defaultGroup := r.Group("/")
	{
		defaultGroup.GET("/", itying.DefaultController{}.Index)
		defaultGroup.GET("/get01", itying.DefaultController{}.Get01)
		defaultGroup.GET("/get02", itying.DefaultController{}.Get02)
		defaultGroup.GET("/delete", itying.DefaultController{}.Delete)
		defaultGroup.GET("/news", itying.DefaultController{}.News)
		defaultGroup.GET("/about", itying.DefaultController{}.About)
	}
}
