package models

type Account struct {
	ID     int `gorm:"primaryKey"`
	Name   string
	Money  int
	Secret Secret
}

type Secret struct {
	ID        int `gorm:"primaryKey"`
	AccountID int
	Password  string
}

func (Secret) TableName() string {
	return "secret"
}

func (Account) TableName() string {
	return "bank"
}
