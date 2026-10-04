// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"

	"rpc/gateway/internal/svc"
	"rpc/gateway/internal/types"
	"rpc/userrpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserNameLogic {
	return &UserNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserNameLogic) UserName(req *types.UserNameReq) (resp *types.UserNameResp, err error) {
	r, err := l.svcCtx.User.GetName(l.ctx, &userclient.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.UserNameResp{Name: r.Name}, nil
}
