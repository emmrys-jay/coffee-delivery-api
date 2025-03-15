package models

import (
	"time"
)

const (
	PAYMENT_PENDING    = "PENDING"
	PAYMENT_PROCESSING = "PROCESSING"
	PAYMENT_COMPLETED  = "COMPLETED"
	PAYMENT_FAILED     = "FAILED"

	ORDER_STATUS_PENDING   = "PENDING"
	ORDER_STATUS_COMPLETED = "COMPLETED"
	ORDER_STATUS_CANCELED  = "CANCELED"
)

func IsValidPaymentStatus(status string) bool {
	return status == PAYMENT_COMPLETED ||
		status == PAYMENT_PROCESSING ||
		status == PAYMENT_PENDING ||
		status == PAYMENT_FAILED
}

func IsValidStatus(status string) bool {
	return status == ORDER_STATUS_CANCELED ||
		status == ORDER_STATUS_COMPLETED ||
		status == ORDER_STATUS_PENDING
}

type Order struct {
	Id          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null" json:"user_id"`
	User        User        `gorm:"not null" json:"user"`
	Status      string      `gorm:"not null" json:"status"`
	TotalAmount string      `gorm:"type:decimal(10,2)" json:"total_amount"`
	OrderItems  []OrderItem `gorm:"not null" json:"order_items,omitempty"`
	CreatedAt   time.Time   `gorm:"not null,index" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"not null" json:"updated_at"`
}

type OrderItem struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"not null,index" json:"order_id"`
	CoffeeID  uint      `json:"coffee_id"`
	Coffee    Coffee    `json:"coffee"`
	Name      string    `json:"name"`
	Quantity  uint      `json:"quantity"`
	UnitPrice string    `gorm:"type:decimal(10,2)" json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemResponse struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"not null,index" json:"order_id"`
	CoffeeID  uint      `json:"coffee_id"`
	Name      string    `json:"name"`
	Quantity  uint      `json:"quantity"`
	UnitPrice string    `gorm:"type:decimal(10,2)" json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
}

type CoffeeInfo struct {
	CoffeeID uint `validate:"required,gte=1" json:"coffee_id"`
	Quantity uint `validate:"required,gte=1" json:"quantity"`
}

type CreateOrderRequest struct {
	Coffees []CoffeeInfo `validate:"required" json:"coffees"`
}

type UpdateOrderRequest struct {
	Status string `json:"status"`
}

type OrderResponse struct {
	Id          uint                `gorm:"primaryKey" json:"id"`
	UserID      uint                `gorm:"not null" json:"user_id"`
	Status      string              `gorm:"not null" json:"status"`
	TotalAmount string              `gorm:"type:decimal(10,2)" json:"total_amount"`
	OrderItems  []OrderItemResponse `gorm:"not null" json:"order_items,omitempty"`
	CreatedAt   time.Time           `gorm:"not null,index" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"not null" json:"updated_at"`
}

func (o *Order) ToOrderResponse() OrderResponse {
	or := OrderResponse{
		Id:          o.Id,
		UserID:      o.UserID,
		Status:      o.Status,
		TotalAmount: o.TotalAmount,
		OrderItems:  []OrderItemResponse{},
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}

	for _, v := range o.OrderItems {
		or.OrderItems = append(or.OrderItems, v.ToOrderItemResponse())
	}

	return or
}

func (o *OrderItem) ToOrderItemResponse() OrderItemResponse {
	return OrderItemResponse{
		Id:        o.Id,
		OrderID:   o.OrderID,
		CoffeeID:  o.CoffeeID,
		Name:      o.Name,
		Quantity:  o.Quantity,
		UnitPrice: o.UnitPrice,
		CreatedAt: o.CreatedAt,
	}
}
