package models

import (
	"time"
)

type StockReservation struct {
	Id               uint      `gorm:"primaryKey" json:"id"`
	CoffeeId         uint      `json:"coffee_id"`
	ReservedQuantity uint      `json:"reserved_quantity"`
	OrderId          uint      `json:"order_id"`
	CreatedAt        time.Time `json:"created_at"`
}
