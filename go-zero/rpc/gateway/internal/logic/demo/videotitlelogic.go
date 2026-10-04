// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"

	"rpc/gateway/internal/svc"
	"rpc/gateway/internal/types"
	"rpc/videorpc/videoclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoTitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoTitleLogic {
	return &VideoTitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoTitleLogic) VideoTitle(req *types.VideoTitleReq) (resp *types.VideoTitleResp, err error) {
	r, err := l.svcCtx.Video.GetTitle(l.ctx, &videoclient.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.VideoTitleResp{Title: r.Title}, nil
}
