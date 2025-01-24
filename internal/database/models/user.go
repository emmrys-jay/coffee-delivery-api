package models

import "time"

type User struct {
	Id        uint   `gorm:"primarykey"`
	FirstName string `gorm:"size:255;not null" validate:"required"`
	LastName  string `gorm:"size:255;not null" validate:"required"`
	Email     string `gorm:"size:255;unique;not null" validate:"required"`
	Password  string `gorm:"size:255;not null" validate:"required"`
	Role      string `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserUpdate struct {
	FirstName string `gorm:"size:255;not null" validate:"required"`
	LastName  string `gorm:"size:255;not null" validate:"required"`
}

type LoginRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}
