package model

import "gorm.io/gorm"

// Tweet struct
type Tweet struct {
	gorm.Model
	Content string `gorm:"not null" json:"content" validate:"required,lte=280"`
	UserID  int    `gorm:"not null" json:"user_id" validate:"required"`
}
