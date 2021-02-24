package constant

import "main/model"

// StatusOptions ...
var StatusOptions = []model.Status{
	{ID: 1, Desc: "待定"},
	{ID: 2, Desc: "失败"},
	{ID: 3, Desc: "成功"},
}

// StatusMap ...
var StatusMap = map[int32]string{
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
