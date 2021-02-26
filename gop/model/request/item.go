package request

// ItemPageReq struct
type ItemPageReq struct {
	PageInfo
	Table    string `form:"table" json:"table" binding:"required"`
	Language string `form:"language" json:"language" binding:"required"`
	Sort     int    `form:"sort" json:"sort"`
	Status   int    `form:"status" json:"status"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// UpdateTextReq struct
type UpdateTextReq struct {
	ID       int    `form:"id" json:"id" binding:"required"`
	Table    string `form:"table" json:"table" binding:"required"`
	Language string `form:"language" json:"language" binding:"required"`
	Text     string `form:"text" json:"text" binding:"required"`
}

// DeleteItemReq struct
type DeleteItemReq struct {
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
