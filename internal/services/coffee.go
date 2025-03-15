package services

import (
	"context"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/repository"
	"github.com/emmrys-jay/coffee-delivery-api/util"
)

type CoffeeService struct {
	repo *repository.CoffeeRepository
}

func NewCoffeeService(repo *repository.CoffeeRepository) *CoffeeService {
	return &CoffeeService{repo: repo}
}

func (s *CoffeeService) CreateCoffee(ctx context.Context, req *models.CreateCoffee) (*models.Coffee, error) {
	coffee := models.Coffee{
		Brand:       req.Brand,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		CreatedAt:   util.CurrentTime(),
		UpdatedAt:   util.CurrentTime(),
	}

	return s.repo.Create(ctx, &coffee)
}

func (s *CoffeeService) GetCoffeeByID(ctx context.Context, id uint) (*models.Coffee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CoffeeService) UpdateCoffee(ctx context.Context, id uint, req *models.UpdateCoffee) error {
	coffee := models.Coffee{
		Id:          id,
		Brand:       req.Brand,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		UpdatedAt:   util.CurrentTime(),
	}

	return s.repo.Update(ctx, &coffee)
}

func (s *CoffeeService) DeleteCoffee(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *CoffeeService) ListCoffees(ctx context.Context) ([]models.Coffee, error) {
	return s.repo.GetAll(ctx)
}
