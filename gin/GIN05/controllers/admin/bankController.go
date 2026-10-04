package admin

import (
	"mygin/GIN05/models"

	"github.com/gin-gonic/gin"
)

type BankController struct {
	BaseController
}

func (con BankController) Index(c *gin.Context) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	account := models.Account{Name: "张三"}
	tx.Find(&account)
	account.Money -= 100
	if err := tx.Model(&models.Account{}).Where("name = ?", "张三").
		Update("money", account.Money).Error; err != nil {
		tx.Rollback()
		con.fail(c)
		return
	}
	account.Name = "李四"
	db.Find(&account)
	account.Money += 100
	if err := tx.Model(&models.Account{}).Where("name=?", "李四").
		Update("money", account.Money).Error; err != nil {
		tx.Rollback()
		con.fail(c)
		return
	}
	tx.Commit()
	con.success(c)
}

func (con BankController) CRUD(c *gin.Context) {
	account := models.Account{
		Name:  "tom",
		Money: 2022,
		Secret: models.Secret{
			Password: "123456",
		},
	}
	db.Create(&account)

	//var account []models.Account
	//db.Preload("Secret", "id=?", "1").Find(&account)
	//db.Joins("Secret").Find(&account, "bank.name=?", "jack")
	// c.JSON(200, account)

	//db.Model(&models.Account{}).Where("name = ?", "jack").Update("money", 2001)
	//db.Model(&models.Secret{}).Where("name = ?", "jack").Update("password", "234567")
	// account := models.Account{
	// 	Name:  "jack",
	// 	Money: 3000,
	// 	Secret: models.Secret{
	// 		Name:     "jack",
	// 		Password: "999999",
	// 	},
	// }
	// db.Updates(&account)

	// var account models.Account
	// db.Find(&account, "name=?", "jack")
	// db.Model(&account).Association("Secret").Replace(&models.Secret{
	// 	Password: "345678",
	// })

	// var account models.Account
	// account.ID = 1
	// db.Model(&account).Association("Secret").Clear()
	// db.Delete(&account)

}
