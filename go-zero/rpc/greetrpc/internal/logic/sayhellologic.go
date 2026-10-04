package logic

import (
	"context"

	"rpc/greetrpc/greeter"
	"rpc/greetrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SayHelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSayHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloLogic {
	return &SayHelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SayHelloLogic) SayHello(in *greeter.SayHelloReq) (*greeter.SayHelloResp, error) {
	return &greeter.SayHelloResp{Message: "Hello " + in.Name}, nil
}
