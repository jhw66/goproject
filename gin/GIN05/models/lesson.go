package models

type Lesson struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Students []Student `gorm:"many2many:lesson_student;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Lesson) TableName() string {
	return "lesson"
}
