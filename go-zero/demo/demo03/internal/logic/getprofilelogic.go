// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"

	"demo03/internal/svc"
	"demo03/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	errorx "github.com/zeromicro/x/errors"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile(req *types.ProfileReq) (resp *types.ProfileResp, err error) {
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	if userId != req.TargetId {
		return nil, errorx.New(403, "无权访问")
	}

	r, e := l.svcCtx.UserRepo.FindByID(l.ctx, userId)
	if e != nil {
		return nil, e
	}

	return &types.ProfileResp{
		UserID:    r.ID,
		Username:  r.Username,
		Email:     r.Email,
		CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
