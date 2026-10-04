package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string         `gorm:"column:username;type:varchar(64);not null;default:''"`
	Password  string         `gorm:"column:password;type:varchar(30);not null;default:''"`
	Email     string         `gorm:"column:email;type:varchar(50);"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"colume:deleted_at;index"`
}
