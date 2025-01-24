package repository

import (
	"context"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"gorm.io/gorm"
)

type CoffeeRepository struct {
	db *gorm.DB
}

func NewCoffeeRepository(db *gorm.DB) *CoffeeRepository {
	return &CoffeeRepository{db: db}
}

func (r *CoffeeRepository) Create(ctx context.Context, coffee *models.Coffee) error {
	return r.db.WithContext(ctx).Create(coffee).Error
}

func (r *CoffeeRepository) GetByID(ctx context.Context, id int) (*models.Coffee, error) {
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

func (r *CoffeeRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Coffee{}, id).Error
}
