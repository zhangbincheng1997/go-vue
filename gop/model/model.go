package model

import (
	"time"

	"gorm.io/gorm"
)

// Model struct
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"create_time"`
	UpdatedAt time.Time      `json:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
