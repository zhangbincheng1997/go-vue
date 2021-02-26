package service

import (
	"context"
	"fmt"
	"main/global"
	"mime/multipart"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
)

// Upload ...
func Upload(file *multipart.FileHeader) (string, error) {
	config := global.CONFIG.Qiniu
	putPolicy := storage.PutPolicy{Scope: config.Bucket}
	mac := qbox.NewMac(config.AccessKey, config.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseHTTPS:      config.UseHTTPS,
		UseCdnDomains: config.UseCdnDomains,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openErr := file.Open()
	if openErr != nil {
		global.LOG.Error("上传文件open异常！", zap.Any("err", openErr.Error()))
		return "", openErr
	}
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		global.LOG.Error("上传文件put异常！", zap.Any("err", putErr.Error()))
		return "", openErr
	}
	return global.CONFIG.Qiniu.ImgPath + "/" + ret.Key, nil
}
