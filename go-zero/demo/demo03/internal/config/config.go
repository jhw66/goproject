// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	RefreshSecret string
	RefreshExpire int64
	Mysql         MysqlConf
	Redis         RedisConf
	Cache         CacheConf
}

type MysqlConf struct {
	Dsn             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConf struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

type CacheConf struct {
	UserPrefix string
	UserTTL    int
}
