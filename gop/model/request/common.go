package request

// PageInfo struct
type PageInfo struct {
	Page  int `form:"page" json:"page" binding:"required"`
	Limit int `form:"limit" json:"limit" binding:"required"`
}
