package controller

import (
	"main/global"
	"main/model"
	"main/model/request"
	"main/model/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

// Info ...
func Info(c *gin.Context) {
	user := global.GetAuthUser(c)
	response.OkWithData(c, user)
}

// Register ...
func Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if global.DB.Where("username = ?", req.Username).Take(&model.User{}).Error != nil {
		user := model.User{
			Username:     req.Username,
			Password:     req.Password,
			Role:         "admin",
			Introduction: "I am a super administrator",
			Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			Name:         "Super Admin",
		}
		res := global.DB.Create(&user)
		if err := res.Error; err != nil {
			global.LOG.Error("更新密码失败", zap.Any("err", err))
			response.FailWithData(c, err)
			return
		}
		response.OkWithData(c, res.RowsAffected)
		return
	}
	response.FailWithMessage(c, "用户名已存在")
}

// GetUserList ...
func GetUserList(c *gin.Context) {
	var req request.UserPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	var count int64
	global.DB.Model(model.User{}).Count(&count)
	var list []model.User
	offset := (req.Page - 1) * req.Limit
	if req.Sort {
		global.DB.Order("id desc")
	}
	global.DB.Limit(req.Limit).Offset(offset).Find(&list)
	response.OkWithData(c, gin.H{
		"list":  list,
		"count": count,
	})
}

// DeleteUser ...
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	res := global.DB.Where("username = ?", id).Delete(model.User{})
	if err := res.Error; err != nil {
		global.LOG.Error("删除用户失败", zap.Any("err", err))
		response.FailWithData(c, err)
	} else {
		response.OkWithData(c, res.RowsAffected)
	}
}

// UpdateInfo ...
func UpdateInfo(c *gin.Context) {
	var req request.UpdateInfoReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	user := global.GetAuthUser(c)
	copier.Copy(&user, &req)
	res := global.DB.Model(&user).Updates(&user) // 不会更新空值
	if err := res.Error; err != nil {
		global.LOG.Error("更新信息失败", zap.Any("err", err))
		response.FailWithData(c, err)
	} else {
		global.RDB.Del(user.Username)
		response.OkWithData(c, res.RowsAffected)
	}
}

// UpdatePassword ...
func UpdatePassword(c *gin.Context) {
	var req request.UpdatePasswordReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	user := global.GetAuthUser(c)
	res := global.DB.Model(model.User{}).Where("username = ? and password = ?", user.Username, req.OldPwd).Update("password", req.NewPwd)
	if err := res.Error; err != nil {
		global.LOG.Error("更新密码失败", zap.Any("err", err))
		response.FailWithData(c, err)
	} else {
		response.OkWithData(c, res.RowsAffected)
	}
}
