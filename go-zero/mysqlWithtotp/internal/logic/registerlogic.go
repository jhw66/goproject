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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	username := strings.TrimSpace(req.Username)
	password := req.Password
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	_, err := l.svcCtx.UserInfoModel.FindOneByUsername(l.ctx, username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	res, err := l.svcCtx.UserInfoModel.Insert(l.ctx, &model.UserInfo{
		Username:          username,
		PasswordHash:      passwordHash,
		TotpSecret:        "",
		TotpPendingSecret: "",
		TotpEnabled:       false,
	})
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		Id:       id,
		Username: username,
	}, nil
}
