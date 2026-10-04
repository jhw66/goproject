package db

import (
	"demo03/internal/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(c config.MysqlConf) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: c.Dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("open gorm failed:%w", err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime))
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql failed:%w", err)
	}

	if err := Auto(db); err != nil {
		return nil, fmt.Errorf("autoimgrate failed:%w", err)
	}

	return db, nil
}
