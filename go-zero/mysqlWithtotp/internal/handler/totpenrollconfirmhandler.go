package handler

import (
	"net/http"

	"mysql/internal/logic"
	"mysql/internal/svc"
	"mysql/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func totpEnrollConfirmHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TotpCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		subject, err := accessSubjectFromRequest(r, svcCtx)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTotpEnrollConfirmLogic(r.Context(), svcCtx)
		resp, err := l.TotpEnrollConfirm(subject.UserID, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
