package logic

import (
	"context"
	"errors"
	"strings"

	"mysql/internal/model"
	"mysql/internal/secure"
	"mysql/internal/svc"
	"mysql/internal/twofa"
	"mysql/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TotpEnrollConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpEnrollConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpEnrollConfirmLogic {
	return &TotpEnrollConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpEnrollConfirmLogic) TotpEnrollConfirm(userID int64, req *types.TotpCodeReq) (*types.TotpToggleResp, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, errors.New("code is required")
	}

	user, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if strings.TrimSpace(user.TotpPendingSecret) == "" {
		return nil, errors.New("no pending totp enrollment found")
	}

	secret, err := secure.DecryptString(l.svcCtx.Config.Totp.SecretCipherKey, user.TotpPendingSecret)
	if err != nil {
		return nil, err
	}
	if !twofa.ValidateCode(secret, code) {
		return nil, errors.New("invalid totp code")
	}

	user.TotpSecret = user.TotpPendingSecret
	user.TotpPendingSecret = ""
	user.TotpEnabled = true
	if err := l.svcCtx.UserInfoModel.Update(l.ctx, user); err != nil {
		return nil, err
	}

	return &types.TotpToggleResp{
		Enabled: true,
	}, nil
}
