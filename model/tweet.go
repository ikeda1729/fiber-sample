package model

import "gorm.io/gorm"

// Tweet struct
type Tweet struct {
	gorm.Model
	Content string `gorm:"not null" json:"title"`
}
