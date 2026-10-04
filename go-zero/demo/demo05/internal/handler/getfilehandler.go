package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"demo05/internal/svc"
	"demo05/internal/types"

	"github.com/qiniu/go-sdk/v7/storagev2/downloader"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		key := req.FileURL
		cfg := svcCtx.Config.Qiniu

		var fileURL string

		if cfg.Private {
			// 私有空间：生成带签名、1小时有效的下载 URL
			urlsProvider := downloader.SignURLsProvider(
				downloader.NewStaticDomainBasedURLsProvider([]string{cfg.Domain}),
				downloader.NewCredentialsSigner(svcCtx.QiniuMac),
				&downloader.SignOptions{TTL: time.Hour},
			)

			iter, err := urlsProvider.GetURLsIter(r.Context(), key, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("生成下载链接失败: %v", err), http.StatusInternalServerError)
				return
			}
			var u url.URL
			ok, err := iter.Peek(&u)
			if err != nil || !ok {
				http.Error(w, "生成下载链接失败", http.StatusInternalServerError)
				return
			}
			fileURL = u.String()
		} else {
			// 公开空间：直接拼接访问 URL
			fileURL = cfg.Domain + "/" + key
		}

		http.Redirect(w, r, fileURL, http.StatusFound)
	}
}

// const photosDir = "./photos"

// func GetFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req types.GetFileReq
// 		if err := httpx.Parse(r, &req); err != nil {
// 			httpx.ErrorCtx(r.Context(), w, err)
// 			return
// 		}

// 		key := req.FileURL
// 		cfg := svcCtx.Config.Qiniu
// 		localPath := filepath.Join(photosDir, time.Now().Format("20060102150405")+key)

// 		if err := os.MkdirAll(photosDir, 0755); err != nil {
// 			http.Error(w, fmt.Sprintf("创建目录失败: %v", err), http.StatusInternalServerError)
// 			return
// 		}

// 		if cfg.Private {
// 			urlsProvider := downloader.SignURLsProvider(
// 				downloader.NewStaticDomainBasedURLsProvider([]string{cfg.Domain}),
// 				downloader.NewCredentialsSigner(svcCtx.QiniuMac),
// 				&downloader.SignOptions{TTL: time.Hour},
// 			)
// 			downloadManager := downloader.NewDownloadManager(&downloader.DownloadManagerOptions{})
// 			_, err := downloadManager.DownloadToFile(context.Background(), key, localPath, &downloader.ObjectOptions{
// 				GenerateOptions:      downloader.GenerateOptions{BucketName: cfg.Bucket},
// 				DownloadURLsProvider: urlsProvider,
// 			})
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("下载文件失败: %v", err), http.StatusInternalServerError)
// 				return
// 			}
// 		} else {
// 			fileURL := cfg.Domain + "/" + key
// 			if err := downloadFromURL(fileURL, localPath); err != nil {
// 				http.Error(w, fmt.Sprintf("下载文件失败: %v", err), http.StatusInternalServerError)
// 				return
// 			}
// 		}

// 		httpx.OkJsonCtx(r.Context(), w, map[string]string{
// 			"message": "下载成功",
// 			"path":    localPath,
// 		})
// 	}
// }

// func downloadFromURL(rawURL, dest string) error {
// 	resp, err := http.Get(rawURL)
// 	if err != nil {
// 		return fmt.Errorf("请求失败: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("远程返回状态码 %d", resp.StatusCode)
// 	}

// 	f, err := os.Create(dest)
// 	if err != nil {
// 		return fmt.Errorf("创建文件失败: %w", err)
// 	}
// 	defer f.Close()

// 	if _, err := io.Copy(f, resp.Body); err != nil {
// 		return fmt.Errorf("写入文件失败: %w", err)
// 	}
// 	return nil
// }
