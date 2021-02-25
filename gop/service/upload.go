package service

import (
	"context"
	"errors"
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
	putPolicy := storage.PutPolicy{Scope: global.CONFIG.Qiniu.Bucket}
	mac := qbox.NewMac(global.CONFIG.Qiniu.AccessKey, global.CONFIG.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		UseHTTPS:      global.CONFIG.Qiniu.UseHTTPS,
		UseCdnDomains: global.CONFIG.Qiniu.UseCdnDomains,
		Zone:          &storage.ZoneHuadong,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := file.Open()
	if openError != nil {
		global.LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))

		return "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		global.LOG.Error("function formUploader.Put() Filed", zap.Any("err", putErr.Error()))
		return "", errors.New("function formUploader.Put() Filed, err:" + putErr.Error())
	}
	return global.CONFIG.Qiniu.ImgPath + "/" + ret.Key, nil
}
