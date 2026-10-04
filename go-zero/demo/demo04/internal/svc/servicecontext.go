// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"demo04/internal/config"
	"demo04/internal/middleware"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config              config.Config
	UserAgentMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		UserAgentMiddleware: middleware.NewUserAgentMiddleware().Handle,
	}
}
