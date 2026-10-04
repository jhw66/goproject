package routers

import (
	"fmt"
	api "mygin/GIN05/controllers/api"

	"github.com/gin-gonic/gin"
)

func ApiGroup(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/", func(c *gin.Context) {
			fmt.Println("aaa")
			c.Set("username", "入")
		}, api.ApiController{}.Index, api.ApiController{}.Add2)
		apiGroup.GET("/v1", api.ApiController{}.Add)
		apiGroup.GET("/v2", api.Timelong, api.ApiController{}.Edit)
		apiGroup.GET("/v3", api.Timelong, api.ApiController{}.Add, api.Printaaa, api.ApiController{}.Edit)
	}
}
