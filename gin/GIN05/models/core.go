package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var err error

// 在其他模块引入models时会自动执行init()方法
func init() {
	NewLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "wjhaccount:Wujiahui789@tcp(127.0.0.1:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//SkipDefaultTransaction: true, //禁用事务(不使用时禁用可提高效率)
		Logger: NewLogger,
	})

	// 运行迁移
	err := DB.AutoMigrate(&Account{}, &Secret{})
	if err != nil {
		panic(err)
	}

	// 检查 secret 表是否存在
	if DB.Migrator().HasTable(&Secret{}) {
		fmt.Println("secret 表存在 ✅")
	} else {
		fmt.Println("secret 表不存在 ❌")
	}
	DB.AutoMigrate(&Lesson{}, &Student{})
	DB.SetupJoinTable(&Lesson{}, "student", &LessonStudent{})
	DB.SetupJoinTable(&Student{}, "lesson", &LessonStudent{})
	DB.AutoMigrate(&Lesson{}, &Student{}, &LessonStudent{})
	DB.AutoMigrate(&Girl{}, &Boy{})
	DB.AutoMigrate(&People{})
}
