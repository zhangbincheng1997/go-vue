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

// Login ...
func Login() {
	// pass
}

// Register ...
func Register(req request.RegisterReq) error {
	if !errors.Is(global.DB.Where("username = ?", req.Username).First(&model.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户名存在！")
	}
	user := model.User{Username: req.Username, Password: utils.MD5(req.Password), Role: constant.GUEST}
	return global.DB.Create(&user).Error
}

// GetUserList ...
func GetUserList(req request.UserPageReq) (interface{}, int64, error) {
	db := global.DB.Model(&model.User{})

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var list []model.User
	offset := (req.Page - 1) * req.Limit
	if req.Sort {
		db.Order("id desc")
	}
	err := db.Limit(req.Limit).Offset(offset).Find(&list).Error
	return list, count, err
}

// DeleteUser ...
func DeleteUser(req request.IDReq) error {
	var user model.User
	copier.Copy(&user, &req)
	err := global.DB.Delete(user).Error
	return err
}

// UpdatePassword ...
func UpdatePassword(user *model.User, req request.UpdatePasswordReq) error {
	var u model.User
	if errors.Is(global.DB.Where("username = ? and password = ?", user.Username, utils.MD5(req.OldPwd)).First(&u).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户名不存在或密码错误！")
	}
	err := global.DB.Model(&u).Update("password", utils.MD5(req.NewPwd)).Error
	return err
}

// UpdateInfo ...
func UpdateInfo(user *model.User, req request.UpdateInfoReq) error {
	copier.Copy(&user, &req)
	err := global.DB.Model(&user).Updates(&user).Error // 不会更新空值
	global.RDB.Del(user.Username)                      // 清空缓存
	return err
}
