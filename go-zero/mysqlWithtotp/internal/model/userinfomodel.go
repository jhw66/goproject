package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserInfoModel = (*customUserInfoModel)(nil)

type (
	// UserInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInfoModel.
	UserInfoModel interface {
		userInfoModel
		withSession(session sqlx.Session) UserInfoModel
	}

	customUserInfoModel struct {
		*defaultUserInfoModel
	}
)

// NewUserInfoModel returns a model for the database table.
func NewUserInfoModel(conn sqlx.SqlConn) UserInfoModel {
	return &customUserInfoModel{
		defaultUserInfoModel: newUserInfoModel(conn),
	}
}

func (m *customUserInfoModel) withSession(session sqlx.Session) UserInfoModel {
	return NewUserInfoModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultUserInfoModel) GetCount(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
