// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"demo05/internal/config"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/region"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

type ServiceContext struct {
	Config        config.Config
	QiniuUploader *uploader.UploadManager
	QiniuMac      *credentials.Credentials
}

func NewServiceContext(c config.Config) *ServiceContext {
	mac := credentials.NewCredentials(c.Qiniu.AccessKey, c.Qiniu.SecretKey)

	options := &uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	}
	if c.Qiniu.Region != "" {
		options.Options.Regions = region.GetRegionByID(c.Qiniu.Region, true)
	}

	return &ServiceContext{
		Config:        c,
		QiniuUploader: uploader.NewUploadManager(options),
		QiniuMac:      mac,
	}
}
