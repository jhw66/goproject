// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"strings"

	"mysql/internal/auth"
	"mysql/internal/model"
	"mysql/internal/svc"
	"mysql/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	username := strings.TrimSpace(req.Username)
	password := req.Password
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := l.svcCtx.UserInfoModel.FindOneByUsername(l.ctx, username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if !auth.ComparePassword(user.PasswordHash, password) {
		return nil, errors.New("invalid username or password")
	}

	if user.TotpEnabled {
		challengeToken, _, err := auth.GenerateToken(
			l.svcCtx.Config.Auth.ChallengeSecret,
			l.svcCtx.Config.Auth.ChallengeExpire,
			user.Id,
			user.Username,
			auth.TokenTypePre2FA,
		)
		if err != nil {
			return nil, err
		}

		return &types.LoginResp{
			Need2FA:        true,
			ChallengeToken: challengeToken,
		}, nil
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

	return &types.LoginResp{
		Need2FA:      false,
		AccessToken:  accessToken,
		AccessExpire: accessExpire,
	}, nil
}
