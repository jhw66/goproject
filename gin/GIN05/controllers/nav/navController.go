package nav

import (
	"mygin/GIN05/models"
	modles "mygin/GIN05/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NavController struct{}

var db = modles.DB

// 添加数据
func (con NavController) Add(c *gin.Context) {
	// girl := &modles.Girl{
	// 	Name: "lili",
	// 	BoyList: []modles.Boy{
	// 		{Name: "haohao"},
	// 		{Name: "haha"},
	// 		{Name: "hh"},
	// 	},
	// }

	// boy := &modles.Boy{
	// 	Name: "wuwu",
	// 	Girl: modles.Girl{
	// 		Name: "miaomiao",
	// 	},
	// }

	// err := db.Create(&girl).Error
	// if err == nil {
	// 	c.String(200, "添加数据成功")
	// }
	// err := db.Create(&boy).Error
	// if err == nil {
	// 	c.String(200, "添加数据成功")
	// }

	var girl models.Girl
	girl.Name = "lulu"
	db.Create(&girl)

	db.Create(&modles.Boy{
		Name: "lolo",
		Girl: girl,
	})
}

func (con NavController) Read(c *gin.Context) {
	// 嵌套预加载
	// var girl models.Girl
	// db.Preload("BoyList.Girl").Find(&girl, "name = ?", "lili")
	// c.JSON(200, girl)

	// 带条件的预加载
	// var girl models.Girl
	// db.Preload("BoyList","id = ?",1).Find(&girl, "name = ?", "lili")
	// c.JSON(200, girl)

	//自定义预加载
	var girl models.Girl
	db.Preload("BoyList", func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", []int{1, 2})
	}).Find(&girl, "name = ?", "lili")
	c.JSON(200, girl)

	// var boys []models.Boy
	// db.Model(&models.Girl{ID: 1}).Association("BoyList").Find(&boys)
	// c.JSON(200, boys)

	// count := db.Model(&models.Girl{Name: "lili"}).Association("BoyList").Count()
	// fmt.Println(count)

}
func (con NavController) Update(c *gin.Context) {
	// girl := models.Girl{}
	// db.Take(&girl, "Name = ?", "lulu")
	// db.Create(&models.Boy{
	// 	Name: "lolo",
	// 	Girl: girl,
	// })

	var girl models.Girl
	db.Take(&girl, 4)
	var boy models.Boy
	db.Where("name = ?", "wuwu").Take(&boy)
	boy.Girl = girl
	db.Save(&boy)

	// 	newBoys := []models.Boy{
	// 		{Name: "jj"},
	// 		{Name: "hh"},
	// 	}
	// 	db.Model(&models.Girl{ID: 2}).Association("BoyList").Replace(&newBoys)

	// 	db.Model(&models.Girl{ID: 3}).Association("BoyList").Append([]models.Boy{{Name: "huahua"}, {Name: "chongchong"}})

}

func (con NavController) Delete(c *gin.Context) {
	// var girl models.Girl
	// girl.ID = 4
	// db.Delete(&girl)

	// db.Preload("BoyList").Find(&girl)
	// db.Model(&girl).Association("BoyList").Delete(&girl.BoyList)

}
