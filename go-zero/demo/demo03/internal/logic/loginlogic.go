// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"demo03/internal/model"
	"demo03/internal/svc"
	"demo03/internal/types"

	"github.com/golang-jwt/jwt/v5"
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

func generateAccessToken(secret string, exp int64, userId int64) (string, error) {
	now := time.Now()
	claim := jwt.MapClaims{
		"userId": userId,
		"iat":    now.Unix(),
		"exp":    jwt.NewNumericDate(now.Add(time.Duration(exp) * time.Second)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
}

func generateRefreshToken(secret string, exp int64, userId int64) (string, error) {
	now := time.Now()
	claim := jwt.MapClaims{
		"userId": userId,
		"iat":    now.Unix(),
		"exp":    jwt.NewNumericDate(now.Add(time.Duration(exp) * time.Second)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret))
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}
	user, err = l.svcCtx.UserRepo.Create(l.ctx, user)
	if err != nil {
		return nil, err
	}
	// if !checkPassword(req.Password, user.Password) {
	// 	return nil, errorx.New(401, "用户名或密码错误")
	// }
	accessToken, _ := generateAccessToken(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, user.ID)
	refreshToken, _ := generateRefreshToken(l.svcCtx.Config.RefreshSecret, l.svcCtx.Config.RefreshExpire, user.ID)

	return &types.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second),
	}, nil
}
