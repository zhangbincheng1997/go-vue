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
	NewPwd string `form:"newPwd" json:"newPwd" binding:"required,min=3,max=12"`
	OldPwd string `form:"oldPwd" json:"oldPwd" binding:"required,min=3,max=12"`
}

// UpdateInfoReq struct
type UpdateInfoReq struct {
	Introduction string `form:"introduction" json:"introduction" binding:"required"`
	Avatar       string `form:"avatar" json:"avatar" binding:"required"`
	Name         string `form:"name" json:"name" binding:"required"`
}

// UpdateRoleReq struct
type UpdateRoleReq struct {
	ID   int    `form:"id" json:"id" binding:"required"`
	Role string `form:"role" json:"role" binding:"required"`
}
