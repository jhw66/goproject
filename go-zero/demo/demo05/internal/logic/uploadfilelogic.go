// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"demo05/internal/svc"
	"demo05/internal/types"

	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(file multipart.File, header *multipart.FileHeader) (*types.UploadResp, error) {
	const maxSize = 10 << 20 // 10 MB
	if header.Size > maxSize {
		return nil, errors.New("文件超过 10 MB 限制")
	}

	// 生成唯一对象名，防止覆盖
	key := fmt.Sprintf("%d_%s", time.Now().UnixNano(),
		filepath.Base(filepath.Clean(header.Filename)))

	err := l.svcCtx.QiniuUploader.UploadReader(l.ctx, file, &uploader.ObjectOptions{
		BucketName: l.svcCtx.Config.Qiniu.Bucket,
		ObjectName: &key,
		FileName:   header.Filename,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("上传到七牛云失败: %w", err)
	}

	fileURL := l.svcCtx.Config.Qiniu.Domain + "/" + key

	return &types.UploadResp{
		Filename: header.Filename,
		Size:     header.Size,
		URL:      fileURL,
	}, nil
}
