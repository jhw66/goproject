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

func UserRoleDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRoleDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewUserRoleDeleteLogic(r.Context(), svcCtx)
		resp, err := l.UserRoleDelete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
