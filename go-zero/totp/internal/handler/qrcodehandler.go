package handler

import (
	"net/http"

	"totp/internal/logic"
	"totp/internal/svc"
	"totp/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// qrcodePNGGetHandler：GET /qrcode?url=<url 编码后的 otpauthUrl>，适合静态 <img src="...">。
func qrcodePNGGetHandler(*svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		otpauthURL := r.URL.Query().Get("url")
		png, err := logic.OTPAuthURLQRCodePNG(otpauthURL, 1024)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(png)
	}
}

// qrcodePNGPostHandler：POST JSONBody，避免 otpauth 出现在 URL 与访问日志的 query 里（更稳妥）。
func qrcodePNGPostHandler(*svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QRCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		png, err := logic.OTPAuthURLQRCodePNG(req.OTPAuthURL, 256)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(png)
	}
}
