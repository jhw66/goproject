package db

import (
	"demo03/internal/model"

	"gorm.io/gorm"
)

func Auto(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}
