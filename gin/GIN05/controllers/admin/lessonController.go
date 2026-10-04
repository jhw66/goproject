package admin

// import (
// 	models "mygin/GIN05/models"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type Lesson struct{}

// func (Lesson) Index(c *gin.Context) {
// 	LessonList := []models.Lesson{}
// 	// db.Preload("LessonKey").Find(&LessonList)
// 	// db.Preload("LessonKey").Where("name=?", "计算机网络").Find(&LessonList)
// 	// db.Preload("LessonKey", "id not in(?,?)", 1, 2).Find(&LessonList)
// 	db.Preload("LessonKey", func(db *gorm.DB) *gorm.DB {
// 		return db.Order("id DESC").Where("id=1")
// 	}).Where("id=1").Find(&LessonList)
// 	c.JSON(200, gin.H{
// 		"result": LessonList,
// 	})
// }
