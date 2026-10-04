package handler

import (
	"net/http"

	"mysql/internal/logic"
	"mysql/internal/svc"
	"mysql/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func totpVerifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TotpVerifyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTotpVerifyLogic(r.Context(), svcCtx)
		resp, err := l.TotpVerify(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
