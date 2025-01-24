package models

import "time"

type Coffee struct {
	Id          int     `gorm:"primaryKey"`
	Brand       string  `validate:"required"`
	Name        string  `validate:"required"`
	Description string  `validate:"required"`
	Price       float64 `validate:"required,gt=0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
