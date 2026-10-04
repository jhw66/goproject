package handler

import (
	"net/http"

	"mysql/internal/auth"
	"mysql/internal/svc"
)

func accessSubjectFromRequest(r *http.Request, svcCtx *svc.ServiceContext) (*auth.Subject, error) {
	token, err := auth.BearerToken(r)
	if err != nil {
		return nil, err
	}

	return auth.ParseToken(svcCtx.Config.Auth.AccessSecret, token, auth.TokenTypeAccess)
}
