// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"demo01/internal/svc"
	"demo01/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoUpdateLogic {
	return &UserInfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoUpdateLogic) UserInfoUpdate(req *types.UserInfoUpdateReq) (resp *types.UserInfoUpdateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
