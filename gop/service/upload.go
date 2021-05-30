package service

import (
	"context"
	"fmt"
	"main/global"
	"mime/multipart"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
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
		return "", openErr
	}
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	global.LOG.Infof("文件名：%v", fileKey)
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		return "", putErr
	}
	return global.CONFIG.Qiniu.ImgPath + "/" + ret.Key, nil
}
