package admin

// import (
// 	models "mygin/GIN05/models"

// 	"github.com/gin-gonic/gin"
// )

// type Student struct{}

// func (Student) Index(c *gin.Context) {
// 	studentList := []models.Student{}
// 	// db.Preload("StudentKey").Find(&studentList)
// 	db.Preload("StudentKey").Limit(2).Offset(1).Order("id desc").Find(&studentList)
// 	c.JSON(200, gin.H{
// 		"result": studentList,
// 	})
// }
