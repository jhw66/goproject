// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"demo03/internal/config"
	"demo03/internal/repository"
	"demo03/pkg/cache"
	"demo03/pkg/db"
	"log"

	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config   config.Config
	UserRepo repository.UserRepository
	Redis    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	gormDB, err := db.NewGorm(c.Mysql)
	if err != nil {
		log.Fatalf("init mysql failed:%v", err)
	}

	rdb, err := cache.NewRedis(c.Redis)
	if err != nil {
		log.Fatalf("init redis failed:%v", err)
	}

	return &ServiceContext{
		Config:   c,
		UserRepo: repository.NewUserRepository(gormDB),
		Redis:    rdb,
	}
}
