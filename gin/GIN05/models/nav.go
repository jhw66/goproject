package models

type Girl struct {
	ID      int64
	Name    string `gorm:"size:32"`
	BoyList []Boy  `gorm:"foreignkey:GirlID"`
}

type Boy struct {
	ID     int64
	GirlID int64
	Name   string `gorm:"size:32"`
	Girl   Girl   `gorm:"foreignkey:GirlID"`
}

func (Girl) TableName() string {
	return "girl"
}

func (Boy) TableName() string {
	return "boy"
}
