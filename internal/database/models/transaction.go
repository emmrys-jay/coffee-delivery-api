package models

import "time"

type Transaction struct {
	ID               uint      `json:"id"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	OrderID          uint      `json:"order_id"`
	Order            Order     `json:"-"`
	Reference        string    `gorm:"not null" json:"reference"`
	PaymentID        string    `gorm:"not null" json:"-"`
	PaymentReference string    `gorm:"not null" json:"payment_reference"`
	PaymentStatus    string    `gorm:"not null" json:"payment_status"`
	TotalAmount      string    `gorm:"type:decimal(10,2)" json:"total_amount"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	AuthorizationURL string    `gorm:"-" json:"authorization_url"`
}

type TransactionRequest struct {
	OrderID uint `validate:"required" json:"order_id"`
}
