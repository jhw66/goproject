// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"demo01/internal/svc"
	"demo01/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleUpdateLogic {
	return &UserRoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRoleUpdateLogic) UserRoleUpdate(req *types.UserRoleUpdateReq) (resp *types.UserRoleUpdateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
