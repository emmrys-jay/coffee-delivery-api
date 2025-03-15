package models

import "time"

type User struct {
	Id        uint      `gorm:"primarykey" json:"id"`
	FirstName string    `gorm:"size:255;not null" validate:"required" json:"first_name"`
	LastName  string    `gorm:"size:255;not null" validate:"required" json:"last_name"`
	Email     string    `gorm:"size:255;unique;not null" validate:"required" json:"email"`
	Password  string    `gorm:"size:255;not null" validate:"required" json:"-"`
	Role      string    `gorm:"not null" validate:"required" json:"role"`
	CreatedAt time.Time `gorm:"not null,index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUser struct {
	FirstName string `validate:"required" json:"first_name"`
	LastName  string `validate:"required" json:"last_name"`
	Email     string `validate:"required" json:"email"`
	Password  string `validate:"required" json:"password"`
	Role      string `validate:"required" json:"role"`
}

type UserUpdate struct {
	FirstName string `validate:"required" json:"first_name"`
	LastName  string `validate:"required" json:"last_name"`
}

type LoginRequest struct {
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
}
