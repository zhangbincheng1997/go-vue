package controller

import (
	"main/constant"
	"main/global"
	"main/model/request"
	"main/model/response"
	"main/service"

	"github.com/gin-gonic/gin"
)

// GetRoleOptions ...
// @Tags User
// @Summary 获取角色选项
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} []model.Option
// @Router /v1/user/role [get]
func GetRoleOptions(c *gin.Context) {
	response.OkWithData(c, constant.RoleOptions)
}

// GetUserInfo ...
// @Tags User
// @Summary 获取信息
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	username := global.GetAuthUser(c)
	user, err := service.GetUserInfo(username)
	if err != nil {
		global.LOG.Errorf("获取信息失败：%v", err)
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithData(c, user)
}

// Login ...
// @Tags User
// @Summary 登录
// @Accept json
// @Produce json
// @Param data body request.LoginReq true "LoginReq"
// @Success 200 {object} response.Response
// @Router /v1/user/login [post]
func Login(c *gin.Context) {
	service.Login() // pass
	response.OkWithMsg(c, "登录成功！")
}

// Register ...
// @Tags User
// @Summary 注册
// @Accept json
// @Produce json
// @Param data body request.RegisterReq true "RegisterReq"
// @Success 200 {object} response.Response
// @Router /v1/user/register [post]
func Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.Register(req)
	if err != nil {
		global.LOG.Errorf("注册失败：%v", err)
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithMsg(c, "注册成功！")
}

// GetUserList ...
// @Tags User
// @Summary 获取用户列表
// @Security ApiKeyAuth
// @Produce json
// @Param data query request.UserPageReq true "UserPageReq"
// @Success 200 {object} response.Response
// @Router /v1/user/list [get]
func GetUserList(c *gin.Context) {
	var req request.UserPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	list, count, err := service.GetUserList(req)
	if err != nil {
		global.LOG.Errorf("获取用户列表失败：%v", err)
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithData(c, gin.H{
		"list":  list,
		"count": count,
	})
}

// DeleteUser ...
// @Tags User
// @Summary 删除用户
// @Security ApiKeyAuth
// @Produce json
// @Param data body request.IDReq true "IDReq"
// @Success 200 {object} response.Response
// @Router /v1/user [delete]
func DeleteUser(c *gin.Context) {
	var req request.IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.DeleteUser(req)
	if err != nil {
		global.LOG.Errorf("删除用户失败：%v", err)
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithMsg(c, "删除用户成功！")
}

// UpdatePassword ...
// @Tags User
// @Summary 更新密码
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body request.UpdatePasswordReq true "UpdatePasswordReq"
// @Success 200 {object} response.Response
// @Router /v1/user/password [put]
func UpdatePassword(c *gin.Context) {
	username := global.GetAuthUser(c)
	var req request.UpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.UpdatePassword(username, req)
	if err != nil {
		global.LOG.Errorf("更新密码失败：%v", err)
		response.FailWithMsg(c, err.Error())
		return
	}
	response.OkWithMsg(c, "更新密码成功！")
}

// UpdateInfo ...
// @Tags User
// @Summary 更新信息
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body request.UpdateInfoReq true "UpdateInfoReq"
// @Success 200 {object} response.Response
// @Router /v1/user/info [put]
func UpdateInfo(c *gin.Context) {
	username := global.GetAuthUser(c)
	var req request.UpdateInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.UpdateInfo(username, req)
	if err != nil {
		global.LOG.Errorf("更新信息失败：%v", err)
		response.FailWithMsg(c, err.Error())
	}
	response.OkWithMsg(c, "更新信息成功！")
}

// UpdateRole ...
// @Tags User
// @Summary 更新角色
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body request.UpdateRoleReq true "UpdateRoleReq"
// @Success 200 {object} response.Response
// @Router /v1/user/role [put]
func UpdateRole(c *gin.Context) {
	var req request.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, err.Error())
		return
	}
	err := service.UpdateRole(req)
	if err != nil {
		global.LOG.Errorf("更新角色失败：%v", err)
		response.FailWithMsg(c, err.Error())
	}
	response.OkWithMsg(c, "更新角色成功！")
}
