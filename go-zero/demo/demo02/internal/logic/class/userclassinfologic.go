// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package class

import (
	"context"

	"demo02/internal/svc"
	"demo02/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserClassInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserClassInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserClassInfoLogic {
	return &UserClassInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserClassInfoLogic) UserClassInfo(req *types.UserClassReq) (resp *types.UserClassResp, err error) {
	// todo: add your logic here and delete this line

	return
}
