// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"demo03/internal/svc"
	"demo03/internal/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
	errorx "github.com/zeromicro/x/errors"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshReq) (resp *types.LoginResp, err error) {
	token, err := jwt.ParseWithClaims(req.RefreshToken, jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(l.svcCtx.Config.RefreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errorx.New(401, "Refresh Token 已失效或非法")
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := int64(claims["userId"].(float64))

	user, _ := l.svcCtx.UserRepo.FindByID(l.ctx, userId)
	if user == nil {
		return nil, errorx.New(404, "该用户不存在")
	}
	newAccess, err := generateAccessToken(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, userId)
	if err != nil {
		return nil, err
	}
	newRefresh, err := generateRefreshToken(l.svcCtx.Config.RefreshSecret, l.svcCtx.Config.RefreshExpire, userId)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		ExpiresIn:    time.Now().Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second),
	}, nil

}
