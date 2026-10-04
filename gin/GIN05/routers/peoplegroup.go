package routers

import (
	p "mygin/GIN05/controllers/people"

	"github.com/gin-gonic/gin"
)

func PeopleGroup(r *gin.Engine) {
	navgroup := r.Group("/people")
	{
		navgroup.POST("/create", p.PeopleController{}.Create)
		navgroup.POST("/findone", p.PeopleController{}.FindOne)
		navgroup.POST("/findmany", p.PeopleController{}.FindMany)
		navgroup.POST("/update", p.PeopleController{}.Update)
		navgroup.POST("/delete", p.PeopleController{}.Delete)
		navgroup.POST("/find", p.PeopleController{}.Find)
	}
}
