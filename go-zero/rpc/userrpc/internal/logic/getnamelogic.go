package logic

import (
	"context"

	"rpc/userrpc/internal/svc"
	"rpc/userrpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNameLogic {
	return &GetNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNameLogic) GetName(in *user.IdRequest) (*user.NameResponse, error) {
	names := map[string]string{
		"1": "张三",
		"2": "李四",
		"3": "王五",
	}
	name, ok := names[in.Id]
	if !ok {
		name = "访客(" + in.Id + ")"
	}
	return &user.NameResponse{Name: name}, nil
}
