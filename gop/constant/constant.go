package constant

import "main/model"

// StatusOptions ...
var StatusOptions = []model.Option{
	{ID: 1, Desc: "待定"},
	{ID: 2, Desc: "失败"},
	{ID: 3, Desc: "成功"},
}

// Status const
const (
	_ = iota
	WAITING
	FAILURE
	SUCCESS
)

// RoleOptions ...
var RoleOptions = []model.Option{
	{ID: "admin", Desc: "管理员"},
	{ID: "guest", Desc: "游客"},
}

// Role const
const (
	ADMIN = "admin"
	GUEST = "guest"
)
