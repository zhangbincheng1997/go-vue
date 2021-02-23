package model

// UserRole struct
type UserRole struct {
	ID     uint `gorm:"primarykey"`
	UserID uint `gorm:"column:userId"`
	RoleID uint `gorm:"column:roleId"`
}
