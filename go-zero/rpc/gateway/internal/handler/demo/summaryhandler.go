// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"rpc/gateway/internal/logic/demo"
	"rpc/gateway/internal/svc"
	"rpc/gateway/internal/types"
)

func SummaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SummaryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := demo.NewSummaryLogic(r.Context(), svcCtx)
		resp, err := l.Summary(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
