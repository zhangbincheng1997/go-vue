package main

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

// Model struct
type Model struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"createTime"`
	UpdatedAt time.Time    `json:"updateTime"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
}

// User struct
type User struct {
	Model
	Username     string `gorm:"column:username" json:"-"`
	Password     string `gorm:"column:password" json:"-"`
	Role         string `gorm:"column:role" json:"role"`
	Introduction string `gorm:"column:introduction" json:"introduction"`
	Avatar       string `gorm:"column:avatar" json:"avatar"`
	Name         string `gorm:"column:name" json:"name"`
}

// Role struct
type Role struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"column:name"`
}

// UserRole struct
type UserRole struct {
	ID     uint `gorm:"primarykey"`
	UserID uint `gorm:"column:userId"`
	RoleID uint `gorm:"column:roleId"`
}

// LoginReq struct
type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// RegisterReq struct
type RegisterReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// UserPageReq struct
type UserPageReq struct {
	Page  int    `form:"page" json:"page" binding:"required"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Sort  bool   `form:"sort" json:"sort"`
	Role  string `form:"role" json:"role"`
}

// UpdateInfoReq struct
type UpdateInfoReq struct {
	Role         string `form:"role" json:"role"`
	Introduction string `form:"introduction" json:"introduction"`
	Avatar       string `form:"avatar" json:"avatar"`
	Name         string `form:"name" json:"name"`
}

// UpdatePasswordReq struct
type UpdatePasswordReq struct {
	NewPwd string `form:"newPwd" json:"newPwd" binding:"required"`
	OldPwd string `form:"oldPwd" json:"oldPwd" binding:"required"`
}

// Generator struct
type Generator struct {
	ID         int    `bson:"id"`
	Collection string `bson:"collection"`
}

// Result struct
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Status struct
type Status struct {
	ID   int    `json:"id"`
	Desc string `json:"desc"`
}

// PageReq struct
type PageReq struct {
	Table    string `form:"table" json:"table" binding:"required"`
	Language string `form:"language" json:"language" binding:"required"`
	Page     int64  `form:"page" json:"page" binding:"required"`
	Limit    int64  `form:"limit" json:"limit" binding:"required"`
	Sort     int    `form:"sort" json:"sort"`
	Status   int    `form:"status" json:"status"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// RecordReq struct
type RecordReq struct {
	Table    string `form:"table" json:"table" binding:"required"`
	Language string `form:"language" json:"language" binding:"required"`
	ID       int    `form:"id" json:"id"`
}

// UpdateTextReq struct
type UpdateTextReq struct {
	ID       int    `form:"id" json:"id" binding:"required"`
	Table    string `form:"table" json:"table" binding:"required"`
	Language string `form:"language" json:"language" binding:"required"`
	Text     string `form:"text" json:"text" binding:"required"`
}

// DeleteReq struct
type DeleteReq struct {
	Table string `form:"table" json:"table" binding:"required"`
	Ids   []int  `form:"ids" json:"ids" binding:"required"`
}

// StatusReq struct
type StatusReq struct {
	Table    string `form:"table" json:"table" binding:"required"`
	Language string `form:"language" json:"language" binding:"required"`
	Ids      []int  `form:"ids" json:"ids" binding:"required"`
	Status   int    `form:"status" json:"status" binding:"required"`
}

var statusOptions = []Status{
	{1, "待定"},
	{2, "失败"},
	{3, "成功"},
}

var statusMap = map[int32]string{
	1: "待定",
	2: "失败",
	3: "成功",
}

// Status const
const (
	_ = iota
	WAITING
	FAILURE
	SUCCESS
)

var identityKey = "username"

func getAuthUser(c *gin.Context) *User {
	user, _ := c.Get(identityKey)
	return user.(*User)
}
