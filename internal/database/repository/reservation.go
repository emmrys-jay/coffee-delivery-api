package repository

import (
	"context"
	"fmt"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"gorm.io/gorm"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (rr *ReservationRepository) ReserveProducts(ctx context.Context, reqs []models.StockReservation) error {
	if err := rr.db.WithContext(ctx).Create(&reqs).Error; err != nil {
		return fmt.Errorf("error creating reservations: %w", err)
	}

	return nil
}

func (rr *ReservationRepository) FetchReservations(ctx context.Context, coffeeId uint) ([]models.StockReservation, error) {
	var reservations []models.StockReservation
	if err := rr.db.WithContext(ctx).Where("coffee_id = ?", coffeeId).Find(&reservations).Error; err != nil {
		return nil, fmt.Errorf("error fetching reservations: %w", err)
	}

	return reservations, nil
}

func (rr *ReservationRepository) CountReservationQuantity(ctx context.Context, coffeeId uint) (int64, error) {
	var result struct {
		TotalQuantity int64
	}
	if err := rr.db.WithContext(ctx).Model(&models.StockReservation{}).Where("coffee_id = ?", coffeeId).Select("SUM(quantity) as total_quantity").Scan(&result).Error; err != nil {
		return -1, fmt.Errorf("error counting reservation quantity: %w", err)
	}

	return result.TotalQuantity, nil
}

func (rr *ReservationRepository) DeleteOrderReservations(ctx context.Context, orderId uint) error {
	if err := rr.db.WithContext(ctx).Where("order_id = ?", orderId).Delete(&models.StockReservation{}).Error; err != nil {
		return fmt.Errorf("error deleting reservations: %w", err)
	}

	return nil
}
