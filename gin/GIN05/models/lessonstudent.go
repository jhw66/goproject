package models

import (
	"time"

	"gorm.io/gorm"
)

type LessonStudent struct {
	LessonID  int `gorm:"primaryKey"`
	StudentID int `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	Student   Student `gorm:"foreignKey:StudentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Lesson    Lesson  `gorm:"foreignKey:LessonID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (LessonStudent) TableName() string {
	return "lesson_student"
}
