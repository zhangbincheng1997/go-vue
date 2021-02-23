package model

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
