// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"io"
	"net/http"

	"demo05/internal/logic"
	"demo05/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 最多 32 MB 不落盘，超出部分自动落盘
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.Error(w, err)
			return
		}
		defer file.Close()

		// 通过读取前 512 字节检测 MIME 类型
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		mimeType := http.DetectContentType(buf[:n])
		if false { //!isAllowedType(mimeType)
			http.Error(w, fmt.Sprintf("不支持的文件类型: %s", mimeType), http.StatusUnsupportedMediaType)
			return
		}
		file.Seek(0, io.SeekStart)

		l := logic.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(file, header)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}
