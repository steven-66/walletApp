package storage

import (
	"context"
	"walletApp/model"
	"walletApp/storage/mocks"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	DB *gorm.DB
}

// NewTransactionRepository creates a new instance of TransactionRepositoryImpl
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{DB: db}
}

// NewMockTransactionRepository creates a new instance of TransactionRepository with mocked methods
func NewMockTransactionRepository(doMocks ...func(mock *mock.Mock)) TransactionRepository {
	mockRepo := &mocks.TransactionRepository{}
	for _, mockFunc := range doMocks {
		mockFunc(&mockRepo.Mock)
	}
	return mockRepo
}

// CreateTransaction logs a new transaction in the database
func (r *TransactionRepositoryImpl) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	return r.DB.WithContext(ctx).Create(transaction).Error
}

// GetTransactionsByUserID retrieves all transactions for a specific user
func (r *TransactionRepositoryImpl) GetTransactionsByUserID(ctx context.Context, userID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("timestamp DESC").Find(&transactions).Error
	return transactions, err
}
