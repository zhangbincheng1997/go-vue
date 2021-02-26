package request

// LoginReq struct
type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=12"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=12"`
}

// RegisterReq struct
type RegisterReq struct {
	Username string `form:"username" json:"username" binding:"required,min=3,max=12"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=12"`
}

// UserPageReq struct
type UserPageReq struct {
	Page  int    `form:"page" json:"page" binding:"required"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Sort  bool   `form:"sort" json:"sort"`
	Role  string `form:"role" json:"role"`
}

// UpdatePasswordReq struct
type UpdatePasswordReq struct {
	NewPwd string `form:"newPwd" json:"newPwd" binding:"required"`
	OldPwd string `form:"oldPwd" json:"oldPwd" binding:"required"`
}

// UpdateInfoReq struct
type UpdateInfoReq struct {
	Introduction string `form:"introduction" json:"introduction"`
	Avatar       string `form:"avatar" json:"avatar"`
	Name         string `form:"name" json:"name"`
}

// UpdateRoleReq struct
type UpdateRoleReq struct {
	ID   int    `form:"id" json:"id"`
	Role string `form:"role" json:"role"`
}
