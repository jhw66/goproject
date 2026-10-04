package peopleController

import (
	"encoding/json"
	"errors"
	"mygin/GIN05/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PeopleController struct{}

var db = models.DB

func (peo PeopleController) Create(c *gin.Context) {
	var people []models.People
	for i := 0; i < 5; i++ {
		people = append(people, models.People{
			//Name: fmt.Sprintf("wjh%d", i),
			Sex: true,
			Age: 11 + i,
		})
	}
	err := db.Create(&people).Error
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.String(200, "增加数据成功")
}

func (peo PeopleController) FindOne(c *gin.Context) {
	var people [3]models.People
	// people[0].ID = 4
	err := db.Take(&people[0], "id= ?", 6).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, "Take "+err.Error())
			return
		} else {
			c.JSON(500, "sql错误")
		}
	}

	err = db.Order("sex ASC").First(&people[1]).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, "Take "+err.Error())
			return
		} else {
			c.JSON(500, "sql错误")
		}
	}

	err = db.Last(&people[2]).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, "Take "+err.Error())
			return
		} else {
			c.JSON(500, "sql错误")
		}
	}

	c.JSON(200, people)
}

func (peo PeopleController) FindMany(c *gin.Context) {
	var peoples []models.People
	// people[0].ID = 4
	tx := db.Order("ID ASC").Find(&peoples, []int{1, 2, 3, 4, 5})
	if tx.Error != nil {
		c.JSON(500, gin.H{"msg": tx.Error.Error()})
		return
	}
	if tx.RowsAffected == 0 {
		c.JSON(404, gin.H{"msg": "没查到数据"})
		return
	}
	datas, _ := json.Marshal(peoples)
	// c.JSON(200, peoples)
	c.Data(200, "application/json", datas)
}

func (pep PeopleController) Update(c *gin.Context) {
	var people models.People
	// err := db.Save(&people).Error
	// if err != nil {
	// 	c.String(500, "更新失败")
	// }
	// c.String(200, "更新成功")
	if err := db.Transaction(func(tx *gorm.DB) (err error) {
		if err := tx.Take(&people, 8).Error; err != nil {
			return err
		}
		people.Name = "SaveName2"
		people.Sex = true
		people.Age = 1
		if err := tx.Save(&people).Error; err != nil {
			return err
		}

		people = models.People{}
		if err := tx.Take(&people, 9).Error; err != nil {
			return err
		}
		people.Name = "SaveName3"
		people.Sex = false
		people.Age = 0
		if err := tx.Select("name").Save(&people).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.String(500, "更新失败")
	} else {
		c.String(200, "更新成功")
	}

	//db.Model(&models.People{}).Where("id<?", 3).Update("name", "UpdateName").Update("age", 18)
	// db.Find(&models.People{}, []int{1, 2, 3}).Updates(models.People{
	// 	Age:  19,
	// 	Name: "UpdatesName",
	// })
}

func (peo PeopleController) Delete(c *gin.Context) {
	var people models.People
	//db.Delete(&people, "name=?", "wjh3")
	db.Where("id<3").Delete(&people)
	//db.Delete(&people)  //拒绝执行 WHERE conditions required
	people.ID = 1
	db.Delete(&people)
}

func (peo PeopleController) Find(c *gin.Context) {
	var peoples []models.People
	// //查询用户名是wjh的
	// db.Where("name = ?", "wjh").Find(&peoples)
	// fmt.Println(peoples)
	// //查询用户名不是wjh的
	// db.Where("name <> ?", "wjh").Find(&peoples)
	// db.Not("name = ?", "wjh").Find(&peoples)
	// fmt.Println(peoples)
	// //查询用户名包括wjh,wjh1的
	// db.Where("name in ?", []string{"wjh", "wjh1"}).Find(&peoples)
	// fmt.Println(peoples)
	// //查询用户名以w开头的
	// db.Where("name like ?", "w%").Find(&peoples)
	// fmt.Println(peoples)
	// //and查询
	// db.Where("age > ? and sex = ?", 23, true).Find(&peoples)
	// fmt.Println(peoples)
	// //or查询
	// //db.Where("age > ? or sex = ?", 23, true).Find(&peoples)
	// db.Or("age > 23").Or("sex = ?", true).Find(&peoples)
	// fmt.Println(peoples)
	// //结构体查询（会过滤零值）
	//db.Where(&models.People{Name: "wjh", Age: 0}).Find(&peoples)
	//fmt.Println(peoples)
	// //map查询（不会过滤零值）
	// db.Where(map[string]interface{}{"name": "wjh", "age": 0}).Find(&peoples)
	// db.Where(map[string]any{"name": "wjh", "age": 0}).Find(&peoples)
	// fmt.Println(peoples)

	// db.Find(&peoples)
	// fmt.Println(peoples)
	// //选择字段访问，其他字段就是零值
	// db.Select("name", "age").Find(&peoples)
	// fmt.Println(peoples)
	// //Scan可以将选择的字段存入另一个结构体中
	// type User struct {
	// 	Name string
	// 	Age  int
	// }
	// var users []User
	// db.Select("name", "age").Find(&peoples).Scan(&users)          //会查询两次
	// db.Model(&models.People{}).Select("name", "age").Scan(&users) //只会查询一次
	// fmt.Println(users)

	// limit := 2
	// page := 1
	// for {
	// 	offset := (page - 1) * limit
	// 	if rows := db.Order("age ASC").Offset(offset).Limit(limit).Find(&peoples).RowsAffected; rows == 0 {
	// 		break
	// 	} else {
	// 		fmt.Printf("第%d页%v\n", page, peoples)
	// 	}
	// 	page++
	// }

	// db.Model(&models.People{}).Distinct("name").Scan(&peoples)
	// db.Model(&models.People{}).Select("distinct name").Scan(&peoples)

	// type Result struct {
	// 	Name  string
	// 	Count int64
	// 	Age   int64
	// }
	// var result []Result
	// db.Model(&peoples).Select("name,COUNT(*) as count,MIN(age) as age").Group("name").Order("age ASC").Scan(&result)
	// sql_mode=only_full_group_by
	// 这个模式在 MySQL 5.7+ 默认开启。
	// 它强制要求：
	// SELECT、ORDER BY 中的非聚合字段必须出现在 GROUP BY 中

	// db.Raw("SELECT * FROM people WHERE id > ?  AND deleted_at IS NULL ORDER BY id DESC", 3).Scan(&peoples)
	// rows := db.Exec("UPDATE people SET age = age + 1 WHERE id = ?  AND deleted_at IS NULL", 4).RowsAffected
	// fmt.Println(rows)

	// db.Select("name", "age", "id").Where("id > (?)", db.Model(&models.People{}).Select("AVG(id)")).Find(&peoples)
	// tx := db.Model(&models.People{}).Select("name,age,id,deleted_at").Where("id > ?", 7) //注意此处一定要把deleted_at提取出来，或者find语句中不包括deleted_at
	// db.Table("(?) as u", tx).Find(&peoples)

	// var p []map[string]interface{}
	// db.Table("people").Where("age > @age and sex = @sex", sql.Named("age", 18), sql.Named("sex", 1)).Find(&p)
	// db.Where("age > @age and sex = @sex", map[string]any{"age": 18, "sex": 1}).Find(&peoples)

	//db.Scopes(WhereAgeB, WhereAgeS).Find(&peoples)

	//c.JSON(200, result)
	//c.JSON(200, p)
	c.JSON(200, peoples)
}

func WhereAgeB(db *gorm.DB) *gorm.DB {
	return db.Where("age > 10")
}

func WhereAgeS(db *gorm.DB) *gorm.DB {
	return db.Where("age < 20")
}
