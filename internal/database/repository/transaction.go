package repository

import (
	"context"
	"errors"
	"time"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	result := r.db.WithContext(ctx).Create(transaction)
	return result.Error
}

func (r *TransactionRepository) UpdateTransactionStatus(ctx context.Context, id int, status string) error {
	result := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction was not found")
	}

	return nil
}

func (r *TransactionRepository) GetTransactionById(ctx context.Context, id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) ListUserTransactions(ctx context.Context, userId uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetPendingTransaction(ctx context.Context, orderId, userId uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(ctx).Where("order_id = ? AND user_id = ? AND payment_status = ?", orderId, userId, models.PAYMENT_PENDING).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}
