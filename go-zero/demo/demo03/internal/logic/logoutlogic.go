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

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) Logout(req *types.LogoutReq) error {
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	if userId != req.TargetId {
		return errorx.New(403, "无权访问")
	}
	return nil
}
