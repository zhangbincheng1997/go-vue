package controller

import (
	"fmt"
	"main/constant"
	"main/global"
	"main/model/request"
	"main/model/response"
	"main/service"
	"path"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetStatus ...
// @Tags Item
// @Summary 获取状态
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} []model.Status
// @Router /v1/item/status [get]
func GetStatus(c *gin.Context) {
	response.OkWithData(c, constant.StatusOptions)
}

// GetList ...
// @Tags Item
// @Summary 获取条目列表
// @Security ApiKeyAuth
// @Produce json
// @Param data query request.ItemPageReq true "ItemPageReq"
// @Success 200 {object} response.Response
// @Router /v1/item/list [get]
func GetList(c *gin.Context) {
	var req request.ItemPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	list, count, err := service.GetItemList(req)
	if err != nil {
		global.LOG.Error("获取条目列表失败！", zap.Any("err", err))
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithData(c, gin.H{
		"list":  list,
		"count": count,
	})
}

// UpdateText ...
// @Tags Item
// @Summary 更新条目
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body request.UpdateTextReq true "UpdateTextReq"
// @Success 200 {object} response.Response
// @Router /v1/item/text [put]
func UpdateText(c *gin.Context) { // 添加翻译
	var req request.UpdateTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.UpdateText(req)
	if err != nil {
		global.LOG.Error("更新条目失败！", zap.Any("err", err))
		response.FailWithData(c, err)
		return
	}
	response.OkWithData(c, "更新条目成功！")
}

// UpdateRecordText ...
// @Tags Item
// @Summary 更新条目翻译
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body request.UpdateTextReq true "UpdateTextReq"
// @Success 200 {object} response.Response
// @Router /v1/item/record/text [put]
func UpdateRecordText(c *gin.Context) { // 添加翻译
	var req request.UpdateTextReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.UpdateRecordText(req)
	if err != nil {
		global.LOG.Error("更新条目翻译失败！", zap.Any("err", err))
		response.FailWithData(c, err)
		return
	}
	response.OkWithData(c, "更新条目翻译成功！")
}

// UpdateStatus ...
// @Tags Item
// @Summary 更新条目状态
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body request.StatusReq true "StatusReq"
// @Success 200 {object} response.Response
// @Router /v1/item/status [put]
func UpdateStatus(c *gin.Context) {
	var req request.StatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.UpdateStatus(req)
	if err != nil {
		global.LOG.Error("更新条目状态失败！", zap.Any("err", err))
		response.FailWithData(c, err)
		return
	}
	response.OkWithData(c, "更新条目状态成功！")
}

// DeleteItem ...
// @Tags Item
// @Summary 删除条目
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body request.DeleteReq true "DeleteReq"
// @Success 200 {object} response.Response
// @Router /v1/item [delete]
func DeleteItem(c *gin.Context) {
	var req request.DeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.DeleteItem(req)
	if err != nil {
		global.LOG.Error("删除条目状态失败！", zap.Any("err", err))
		response.FailWithData(c, err)
		return
	}
	response.OkWithData(c, "删除条目状态成功！")
}

// ImportData ...
// @Tags Item
// @Summary 导入数据
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Param table formData file string "table"
// @Success 200 {object} response.Response
// @Router /v1/item/import [post]
func ImportData(c *gin.Context) {
	file, _ := c.FormFile("file")
	table := c.Request.FormValue("table")
	err := service.ImportData(file, table)
	if err != nil {
		global.LOG.Error("导入数据失败！", zap.Any("err", err))
		response.FailWithData(c, err)
		return
	}
	response.OkWithData(c, "导入数据成功！")
}

// ExportData ...
// @Tags Item
// @Summary 导出
// @Security ApiKeyAuth
// @Produce json
// @Param table query string true "table"
// @Param language query string true "language"
// @Success 200 {object} response.Response
// @Router /v1/item/export [get]
func ExportData(c *gin.Context) {
	table := c.Query("table")
	language := c.Query("language")
	dir := "data"
	_filename := "data.csv"
	filename := path.Join(dir, _filename)

	err := service.ExportData(filename, table, language)
	if err != nil {
		global.LOG.Error("导出数据失败！", zap.Any("err", err))
		response.FailWithData(c, err)
		return
	}

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", _filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filename)
}
