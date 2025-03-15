package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/repository"
	"github.com/emmrys-jay/coffee-delivery-api/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo, jwtSecret: os.Getenv("SECRET")}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUser) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if req.Role != "user" && req.Role != "admin" {
		return nil, errors.New("invalid role entered")
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      req.Role,
		CreatedAt: util.CurrentTime(),
		UpdatedAt: util.CurrentTime(),
	}
	return s.repo.CreateUser(ctx, &user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, updateRequest *models.UserUpdate) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	user.FirstName = updateRequest.FirstName
	user.LastName = updateRequest.LastName

	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": fmt.Sprint(user.Id),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"role":    user.Role,
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.ListUsers(ctx)
}
