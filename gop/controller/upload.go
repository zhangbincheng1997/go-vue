package controller

import (
	"main/global"
	"main/model/response"
	"main/service"

	"github.com/gin-gonic/gin"
)

// UploadFile ...
// @Tags Upload
// @Summary 上传文件
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} response.Response
// @Router /v1/upload [post]
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	url, err := service.Upload(file) // 文件上传后拿到文件路径
	if err != nil {
		global.LOG.Errorf("上传失败：%v", err)
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithData(c, url)
}
