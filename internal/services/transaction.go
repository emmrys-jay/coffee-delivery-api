package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/repository"
	"github.com/emmrys-jay/coffee-delivery-api/internal/platform"
	"github.com/emmrys-jay/coffee-delivery-api/internal/platform/paystack"
	"github.com/emmrys-jay/coffee-delivery-api/util"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionService struct {
	orderRepo   *repository.OrderRepository
	userRepo    *repository.UserRepository
	reserveRepo *repository.ReservationRepository
	trxRepo     *repository.TransactionRepository

	provider        platform.PaymentProvider
	paymentPlatform platform.Provider
}

func NewTransactionService(
	provider string,
	orderRepo *repository.OrderRepository,
	userRepo *repository.UserRepository,
	reserveRepo *repository.ReservationRepository,
	trxRepo *repository.TransactionRepository,
) (*TransactionService, error) {

	ps := &TransactionService{
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		reserveRepo: reserveRepo,
		trxRepo:     trxRepo,
	}

	switch provider {
	case platform.PAYSTACK.String():
		ps.provider = platform.PAYSTACK
		ps.paymentPlatform = paystack.NewPaystackService()
	default:
		return nil, errors.New("invalid payment provider")
	}

	return ps, nil
}

func (ps *TransactionService) Initiate(ctx context.Context, userId uint, req *models.TransactionRequest) (models.Transaction, error) {

	user, err := ps.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error getting user, %w", err)
	}

	order, err := ps.orderRepo.GetOrder(ctx, req.OrderID)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error getting order, %w", err)
	}

	totalAmount, err := decimal.NewFromString(order.TotalAmount)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error parsing total amount, %w", err)
	}

	pendingTransaction, err := ps.trxRepo.GetPendingTransaction(ctx, order.Id, user.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return models.Transaction{}, fmt.Errorf("error getting pending transaction, %w", err)
	}

	if pendingTransaction != nil {
		return *pendingTransaction, nil
	}

	reference := util.GenerateReference()

	// Initiate payment transaction
	resp, err := ps.paymentPlatform.InitiateTransaction(user.Email, totalAmount, reference)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error initiating payment transaction, %w", err)
	}

	trx := models.Transaction{
		OrderID:          order.Id,
		UserID:           user.Id,
		Reference:        reference,
		PaymentID:        resp.Data.AccessCode,
		PaymentReference: resp.Data.Reference,
		PaymentStatus:    models.PAYMENT_PENDING,
		CreatedAt:        util.CurrentTime(),
		UpdatedAt:        util.CurrentTime(),
		TotalAmount:      totalAmount.String(),
		AuthorizationURL: resp.Data.AuthorizationURL,
	}

	if err = ps.trxRepo.CreateTransaction(ctx, &trx); err != nil {
		return models.Transaction{}, fmt.Errorf("error creating transaction, %w", err)
	}

	return trx, nil
}
