// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"demo01/internal/logic"
	"demo01/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserRoleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUserRoleListLogic(r.Context(), svcCtx)
		resp, err := l.UserRoleList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
