// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"

	"rpc/gateway/internal/svc"
	"rpc/gateway/internal/types"
	"rpc/greetrpc/greeterclient"
	"rpc/userrpc/userclient"
	"rpc/videorpc/videoclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type SummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SummaryLogic {
	return &SummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SummaryLogic) Summary(req *types.SummaryReq) (resp *types.SummaryResp, err error) {
	name := req.Name
	if name == "" {
		name = "go-zero"
	}
	uid := req.Uid
	if uid == "" {
		uid = "1"
	}
	vid := req.Vid
	if vid == "" {
		vid = "1"
	}

	hello, err := l.svcCtx.Greeter.SayHello(l.ctx, &greeterclient.SayHelloReq{Name: name})
	if err != nil {
		return nil, err
	}
	un, err := l.svcCtx.User.GetName(l.ctx, &userclient.IdRequest{Id: uid})
	if err != nil {
		return nil, err
	}
	vt, err := l.svcCtx.Video.GetTitle(l.ctx, &videoclient.IdRequest{Id: vid})
	if err != nil {
		return nil, err
	}

	return &types.SummaryResp{
		Greeting: hello.Message,
		UserName: un.Name,
		Title:    vt.Title,
	}, nil
}
