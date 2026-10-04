// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package class

import (
	"context"

	"demo02/internal/svc"
	"demo02/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserClassAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClassAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClassAddLogic {
	return &UserClassAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClassAddLogic) UserClassAdd(req *types.UserClassAddReq) (resp *types.UserClassAddResp, err error) {
	// todo: add your logic here and delete this line

	return
}
