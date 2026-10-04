package models

type Student struct {
	ID      int
	Name    string
	Lessons []Lesson `gorm:"many2many:lesson_student;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Student) TableName() string {
	return "student"
}
