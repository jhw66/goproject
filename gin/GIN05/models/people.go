package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type People struct {
	ID        uint
	Name      string
	Age       int `gorm:"tinyint"`
	Sex       bool
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (People) TableName() string {
	return "people"
}
func (u *People) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Name == "" {
		return errors.New("name 不能为空")
	}
	return nil
}
