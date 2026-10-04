package curd

import (
	"mygin/GIN05/models"

	"github.com/gin-gonic/gin"
)

type CrudController struct{}

var db = models.DB

func (crud CrudController) Add(c *gin.Context) {
	Lessons := []models.Lesson{
		{
			Name: "math",
			Students: []models.Student{
				{Name: "lili"},
				{Name: "huahua"},
				{Name: "mingming"},
			},
		},
		{
			Name: "english",
			Students: []models.Student{
				{Name: "huahua"},
				{Name: "papa"},
				{Name: "wuwu"},
				{Name: "xingxing"},
			},
		},
	}
	db.Create(&Lessons)

	// IDList := []int64{1, 2}
	// var Lessons []models.Lesson
	// db.Find(&Lessons, "id in ?", IDList)
	// db.Create(&[]models.Student{
	// 	{
	// 		Name:    "huihui",
	// 		Lessons: Lessons,
	// 	},
	// 	{
	// 		Name:    "wuwu",
	// 		Lessons: Lessons,
	// 	},
	// 	{
	// 		Name: "jiajia",
	// 		Lessons: []models.Lesson{
	// 			{
	// 				Name: "chinese",
	// 			},
	// 		},
	// 	},
	// })
}

func (crud CrudController) Read(c *gin.Context) {
	// var Students []models.Student
	// db.Preload("Lessons").Where("name = ?", "xingxing").Find(&Students)
	// c.JSON(200, Students)

	// var Lessons []models.Lesson
	// db.Preload("Students").Where("name = ?", "math").Find(&Lessons)
	// c.JSON(200, Lessons)

	// var students []models.Student
	// db.Model(&models.Student{}).Joins("JOIN lesson_student ls ON ls.student_id=student.id").
	// 	Joins("JOIN lesson ON lesson.id=ls.lesson_id").Scan(&students)
	// c.JSON(200, students)

	type Resultstudent struct {
		Name   string
		Lesson string
	}
	var results []Resultstudent
	db.Table("student as s").Select("s.name as name,l.name as lesson").Joins("JOIN lesson_student ls ON ls.student_id=s.id").
		Joins("JOIN lesson l ON l.id=ls.lesson_id").Scan(&results)
	c.JSON(200, results)

}

func (crud CrudController) Update(c *gin.Context) {
	// var lesson models.Lesson
	// db.Where("name = ?", "math").Take(&lesson)
	// db.Model(&lesson).Association("Students").Replace([]models.Student{
	// 	{Name: "mama"},
	// 	{Name: "baba"},
	// })

	var lesson models.Lesson
	db.Where("name = ?", "math").Take(&lesson)
	db.Model(&lesson).Association("Students").Append([]models.Student{
		{ID: 35},
		{ID: 36},
	})

}

func (crud CrudController) Delete(c *gin.Context) {
	// var Students []models.Student
	// db.Find(&Students)
	// db.Delete(&Students)

	var lesson models.Lesson
	db.Preload("Students").Where("name = ?", "english").Take(&lesson)
	db.Model(&lesson).Association("Students").Delete(&lesson.Students)
}
