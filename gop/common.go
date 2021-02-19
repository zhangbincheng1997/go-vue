package main

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

// User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
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
