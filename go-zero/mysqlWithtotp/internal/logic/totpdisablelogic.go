package logic

import (
	"context"
	"errors"
	"strings"

	"mysql/internal/auth"
	"mysql/internal/model"
	"mysql/internal/secure"
	"mysql/internal/svc"
	"mysql/internal/twofa"
	"mysql/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TotpDisableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpDisableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpDisableLogic {
	return &TotpDisableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpDisableLogic) TotpDisable(userID int64, req *types.Disable2FAReq) (*types.TotpToggleResp, error) {
	password := req.Password
	code := strings.TrimSpace(req.Code)
	if password == "" || code == "" {
		return nil, errors.New("password and code are required")
	}

	user, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if !user.TotpEnabled || strings.TrimSpace(user.TotpSecret) == "" {
		return nil, errors.New("2fa is not enabled")
	}
	if !auth.ComparePassword(user.PasswordHash, password) {
		return nil, errors.New("invalid password or totp code")
	}

	secret, err := secure.DecryptString(l.svcCtx.Config.Totp.SecretCipherKey, user.TotpSecret)
	if err != nil {
		return nil, err
	}
	if !twofa.ValidateCode(secret, code) {
		return nil, errors.New("invalid password or totp code")
	}

	user.TotpSecret = ""
	user.TotpPendingSecret = ""
	user.TotpEnabled = false
	if err := l.svcCtx.UserInfoModel.Update(l.ctx, user); err != nil {
		return nil, err
	}

	return &types.TotpToggleResp{
		Enabled: false,
	}, nil
}
