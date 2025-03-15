package repository

import (
	"context"
	"fmt"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) (uint, error) {

	if err := r.db.WithContext(ctx).Create(&order).Error; err != nil {
		return 0, fmt.Errorf("error creating order: %w", err)
	}

	return order.Id, nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	if err := r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Order{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) ListUserOrders(ctx context.Context, id uint) ([]models.Order, error) {
	var order []models.Order
	if err := r.db.WithContext(ctx).Preload("OrderItems").Where("user_id = ?", id).Find(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}
