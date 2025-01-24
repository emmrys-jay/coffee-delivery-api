package services

import (
	"context"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/repository"
)

type CoffeeService struct {
	repo *repository.CoffeeRepository
}

func NewCoffeeService(repo *repository.CoffeeRepository) *CoffeeService {
	return &CoffeeService{repo: repo}
}

func (s *CoffeeService) CreateCoffee(ctx context.Context, coffee *models.Coffee) error {
	return s.repo.Create(ctx, coffee)
}

func (s *CoffeeService) GetCoffeeByID(ctx context.Context, id int) (*models.Coffee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CoffeeService) UpdateCoffee(ctx context.Context, coffee *models.Coffee) error {
	return s.repo.Update(ctx, coffee)
}

func (s *CoffeeService) DeleteCoffee(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *CoffeeService) ListCoffees(ctx context.Context) ([]models.Coffee, error) {
	return s.repo.GetAll(ctx)
}
