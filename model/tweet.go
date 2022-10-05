package model

import "gorm.io/gorm"

// Tweet struct
type Tweet struct {
	gorm.Model
	Content string `gorm:"not null" json:"content" validate:"required,lte=280"`
	UserID  string `gorm:"not null"`
	User    User   `gorm:"reference:Username,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
