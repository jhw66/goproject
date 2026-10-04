package logic

import (
	"context"
	"errors"

	"mysql/internal/model"
	"mysql/internal/secure"
	"mysql/internal/svc"
	"mysql/internal/twofa"
	"mysql/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TotpEnrollStartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTotpEnrollStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotpEnrollStartLogic {
	return &TotpEnrollStartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TotpEnrollStartLogic) TotpEnrollStart(userID int64) (*types.TotpEnrollStartResp, error) {
	user, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if user.TotpEnabled {
		return nil, errors.New("2fa is already enabled")
	}

	enrollment, err := twofa.GenerateEnrollment(
		l.svcCtx.Config.Totp.Issuer,
		user.Username,
		l.svcCtx.Config.Totp.QRCodeSize,
	)
	if err != nil {
		return nil, err
	}

	encryptedSecret, err := secure.EncryptString(l.svcCtx.Config.Totp.SecretCipherKey, enrollment.Secret)
	if err != nil {
		return nil, err
	}

	user.TotpPendingSecret = encryptedSecret
	if err := l.svcCtx.UserInfoModel.Update(l.ctx, user); err != nil {
		return nil, err
	}

	return &types.TotpEnrollStartResp{
		Issuer:        enrollment.Issuer,
		AccountName:   enrollment.AccountName,
		OTPAuthURL:    enrollment.OTPAuthURL,
		QRCodeDataURL: enrollment.QRCodeDataURL,
	}, nil
}
