package model

// Role struct
type Role struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"column:name"`
}
