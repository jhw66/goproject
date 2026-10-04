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

type GenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateLogic {
	return &GenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateLogic) Generate(req *types.GenerateReq) (*types.GenerateResp, error) {
	accountName := strings.TrimSpace(req.AccountName)
	if accountName == "" {
		return nil, errors.New("accountName is required")
	}

	issuer := strings.TrimSpace(req.Issuer)
	if issuer == "" {
		issuer = strings.TrimSpace(l.svcCtx.Config.DefaultIssuer)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return nil, err
	}

	return &types.GenerateResp{
		Secret:      key.Secret(),
		Issuer:      issuer,
		AccountName: accountName,
		OTPAuthURL:  key.URL(),
	}, nil
}
