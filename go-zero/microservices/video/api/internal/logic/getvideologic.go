// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"rpc/types/user"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoLogic {
	return &GetVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVideoLogic) GetVideo(req *types.VideoReq) (resp *types.VideoRes, err error) {
	// todo: add your logic here and delete this line

	userl, err := l.svcCtx.UserRpc.GetUser(l.ctx, &user.IdRequest{
		Id: "1",
	})

	if err != nil {
		return nil, err
	}

	return &types.VideoRes{
		Id:   userl.Id,
		Name: userl.Name,
	}, nil
}
