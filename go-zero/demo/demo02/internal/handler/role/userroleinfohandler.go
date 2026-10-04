// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"net/http"

	"demo02/internal/logic/role"
	"demo02/internal/svc"
	"demo02/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserRoleInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewUserRoleInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserRoleInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
