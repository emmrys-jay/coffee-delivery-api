package repository

import (
	"context"
	"fmt"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"gorm.io/gorm"
)

type CoffeeRepository struct {
	db *gorm.DB
}

func NewCoffeeRepository(db *gorm.DB) *CoffeeRepository {
	return &CoffeeRepository{db: db}
}

func (r *CoffeeRepository) Create(ctx context.Context, coffee *models.Coffee) (*models.Coffee, error) {
	err := r.db.WithContext(ctx).Create(coffee).Error
	return coffee, err
}

func (r *CoffeeRepository) GetByID(ctx context.Context, id uint) (*models.Coffee, error) {
	var coffee models.Coffee
	if err := r.db.WithContext(ctx).First(&coffee, id).Error; err != nil {
		return nil, err
	}
	return &coffee, nil
}

func (r *CoffeeRepository) GetAll(ctx context.Context) ([]models.Coffee, error) {
	var coffees []models.Coffee
	if err := r.db.WithContext(ctx).Find(&coffees).Error; err != nil {
		return nil, err
	}
	return coffees, nil
}

func (r *CoffeeRepository) Update(ctx context.Context, coffee *models.Coffee) error {
	return r.db.WithContext(ctx).Save(coffee).Error
}

func (r *CoffeeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Coffee{}, id).Error
}

func (r *CoffeeRepository) GetMany(ctx context.Context, ids []uint) ([]models.Coffee, error) {
	var coffees []models.Coffee
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&coffees).Error; err != nil {
		return nil, err
	}
	return coffees, nil
}

func (r *CoffeeRepository) SubtractReservations(ctx context.Context, reqs []models.StockReservation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, req := range reqs {
			if err := tx.Model(&models.Coffee{}).Where("id = ?", req.CoffeeId).Update("quantity", gorm.Expr("quantity - ?", req.ReservedQuantity)).Error; err != nil {
				return fmt.Errorf("error subtracting from coffee with id %v: %w", req.CoffeeId, err)
			}
		}
		return nil
	})
}
