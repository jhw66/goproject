// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"

	"rpc/gateway/internal/svc"
	"rpc/gateway/internal/types"
	"rpc/greetrpc/greeterclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelloLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
	return &HelloLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloLogic) Hello(req *types.HelloReq) (resp *types.HelloResp, err error) {
	r, err := l.svcCtx.Greeter.SayHello(l.ctx, &greeterclient.SayHelloReq{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &types.HelloResp{Message: r.Message}, nil
}
