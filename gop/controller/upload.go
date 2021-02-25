package controller

import (
	"main/global"
	"main/model/response"
	"main/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Upload ...
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	url, err := service.Upload(file) // 文件上传后拿到文件路径
	if err != nil {
		global.LOG.Error("上传失败！", zap.Any("err", err))
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithData(c, url)
}
