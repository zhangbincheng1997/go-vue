package service

import (
	"errors"
	"main/constant"
	"main/global"
	"main/model"
	"main/model/request"
	"main/utils"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// GetUserInfo ...
func GetUserInfo(username string) (model.User, error) {
	var user model.User
	err := global.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

// Login ...
func Login() {
	// pass
}

// Register ...
func Register(req request.RegisterReq) error {
	if !errors.Is(global.DB.Where("username = ?", req.Username).First(&model.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户名已存在！")
	}
	user := model.User{Username: req.Username, Password: utils.MD5(req.Password), Role: constant.GUEST}
	err := global.DB.Create(&user).Error
	return err
}

// GetUserList ...
func GetUserList(req request.UserPageReq) (interface{}, int64, error) {
	db := global.DB.Model(&model.User{})

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var list []model.User
	if req.Sort {
		db.Order("id desc")
	}
	if req.Role != "" {
		db.Where("role = ?", req.Role)
	}
	offset := (req.Page - 1) * req.Limit
	err := db.Limit(req.Limit).Offset(offset).Find(&list).Error
	return list, count, err
}

// DeleteUser ...
func DeleteUser(req request.IDReq) error {
	err := global.DB.Delete(&model.User{}, req.ID).Error
	return err
}

// UpdatePassword ...
func UpdatePassword(username string, req request.UpdatePasswordReq) error {
	var user model.User
	if errors.Is(global.DB.Where("username = ? and password = ?", username, utils.MD5(req.OldPwd)).First(&user).Error, gorm.ErrRecordNotFound) {
		return errors.New("密码错误！")
	}
	err := global.DB.Model(&user).Update("password", utils.MD5(req.NewPwd)).Error
	return err
}

// UpdateInfo ...
func UpdateInfo(username string, req request.UpdateInfoReq) error {
	var user model.User
	copier.Copy(&user, &req)
	err := global.DB.Where("username = ?", username).Updates(&user).Error // 不会更新空值
	global.RDB.Del(username)                                              // 清空缓存
	return err
}

// UpdateRole ...
func UpdateRole(req request.UpdateRoleReq) error {
	var user model.User
	var username string
	copier.Copy(&user, &req)
	err := global.DB.Model(&user).Updates(&user).Pluck("username", &username).Error // 不会更新空值
	global.RDB.Del(username)                                                        // 清空缓存
	return err
}
