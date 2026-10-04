// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"mysql/internal/config"
	"mysql/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	UserInfoModel model.UserInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:        c,
		UserInfoModel: model.NewUserInfoModel(conn),
	}
}
