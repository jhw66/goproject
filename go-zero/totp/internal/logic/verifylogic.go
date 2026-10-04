// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"strings"

	"totp/internal/svc"
	"totp/internal/types"

	"github.com/pquerna/otp/totp"
	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyLogic {
	return &VerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyLogic) Verify(req *types.VerifyReq) (*types.VerifyResp, error) {
	secret := strings.TrimSpace(req.Secret)
	code := strings.TrimSpace(req.Code)
	if secret == "" || code == "" {
		return nil, errors.New("secret and code are required")
	}

	valid := totp.Validate(code, secret)

	return &types.VerifyResp{
		Valid: valid,
	}, nil
}
