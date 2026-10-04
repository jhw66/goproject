package logic

import (
	"context"

	"rpc/videorpc/internal/svc"
	"rpc/videorpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTitleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTitleLogic {
	return &GetTitleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTitleLogic) GetTitle(in *video.IdRequest) (*video.TitleResponse, error) {
	titles := map[string]string{
		"1": "入门：go-zero RPC",
		"2": "进阶：etcd 服务发现",
		"3": "实战：微服务联调",
	}
	title, ok := titles[in.Id]
	if !ok {
		title = "未命名视频(" + in.Id + ")"
	}
	return &video.TitleResponse{Title: title}, nil
}
