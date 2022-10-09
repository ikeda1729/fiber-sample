package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username  string  `gorm:"unique_index;not null;unique" json:"username" validate:"required,lte=50,gte=5"`
	Email     string  `gorm:"unique_index;not null;unique" json:"email" validate:"required,email,lte=150"`
	Password  string  `gorm:"not null" json:"password" validate:"required,gte=8"`
	Names     string  `json:"names"`
	Tweets    []Tweet `json:"tweets"`
	Followees []*User `gorm:"many2many:user_followees"`
}

type UserResponse struct {
	gorm.Model
	Username string
}

type IsFollowingResponse struct {
	UserID      int
	FolloweeID  string
	IsFollowing bool
}
