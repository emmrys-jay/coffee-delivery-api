package models

import "time"

type Coffee struct {
	Id          uint      `gorm:"primaryKey" json:"id"`
	Brand       string    `json:"brand"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `gorm:"type:decimal(10,2)" json:"price"`
	Quantity    uint      `json:"quantity"`
	CreatedAt   time.Time `gorm:"not null,index" json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCoffee struct {
	Brand       string `validate:"required" json:"brand"`
	Name        string `validate:"required" json:"name"`
	Description string `validate:"required" json:"description"`
	Price       string `validate:"required,sig" json:"price"`
	Quantity    uint   `validate:"required,min=1" json:"quantity"`
}

type UpdateCoffee struct {
	Brand       string `validate:"required" json:"brand"`
	Name        string `validate:"required" json:"name"`
	Description string `validate:"required" json:"description"`
	Price       string `validate:"required,sig" json:"price"`
	Quantity    uint   `validate:"required,min=1" json:"quantity"`
}
