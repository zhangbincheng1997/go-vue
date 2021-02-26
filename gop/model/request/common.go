package request

// PageInfo struct
type PageInfo struct {
	Page  int `form:"page" json:"page" binding:"required"`
	Limit int `form:"limit" json:"limit" binding:"required"`
}

// IDReq struct
type IDReq struct {
	ID uint `form:"id" json:"id" binding:"required"`
}

// IdsReq struct
type IdsReq struct {
	Ids []int `form:"ids" json:"ids" binding:"required"`
}
