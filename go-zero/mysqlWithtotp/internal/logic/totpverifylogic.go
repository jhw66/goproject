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

type TotpVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpVerifyLogic {
	return &TotpVerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpVerifyLogic) TotpVerify(req *types.TotpVerifyReq) (*types.TotpVerifyResp, error) {
	challengeToken := strings.TrimSpace(req.ChallengeToken)
	code := strings.TrimSpace(req.Code)
	if challengeToken == "" || code == "" {
		return nil, errors.New("challengeToken and code are required")
	}

	subject, err := auth.ParseToken(l.svcCtx.Config.Auth.ChallengeSecret, challengeToken, auth.TokenTypePre2FA)
	if err != nil {
		return nil, errors.New("invalid challenge token")
	}

	user, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, subject.UserID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if !user.TotpEnabled || strings.TrimSpace(user.TotpSecret) == "" {
		return nil, errors.New("2fa is not enabled for this account")
	}

	secret, err := secure.DecryptString(l.svcCtx.Config.Totp.SecretCipherKey, user.TotpSecret)
	if err != nil {
		return nil, err
	}
	if !twofa.ValidateCode(secret, code) {
		return nil, errors.New("invalid totp code")
	}

	accessToken, accessExpire, err := auth.GenerateToken(
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
		user.Id,
		user.Username,
		auth.TokenTypeAccess,
	)
	if err != nil {
		return nil, err
	}

	return &types.TotpVerifyResp{
		AccessToken:  accessToken,
		AccessExpire: accessExpire,
	}, nil
}
