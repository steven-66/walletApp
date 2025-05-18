package storage

import (
	"context"
	"walletApp/model"
	"walletApp/storage/mocks"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type balanceRepositoryImpl struct {
	DB *gorm.DB
}

// NewBalanceRepository creates a new instance of balanceRepositoryImpl
func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepositoryImpl{DB: db}
}

// NewMockBalanceRepository creates a new instance of BalanceRepository with mocked methodsgo get github.com/DATA-DOG/go-sqlmock
func NewMockBalanceRepository(doMocks ...func(mock *mock.Mock)) BalanceRepository {
	mockRepo := &mocks.BalanceRepository{}
	for _, mockFunc := range doMocks {
		mockFunc(&mockRepo.Mock)
	}
	return mockRepo
}

// GetBalance retrieves the user's balance from Redis or the database
func (r *balanceRepositoryImpl) GetBalance(ctx context.Context, userID uint) (float64, error) {
	var balance model.Balance
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Select("balance").First(&balance).Error
	if err != nil {
		return 0, err
	}
	return balance.Balance, nil
}

// UpdateBalance updates the user's balance in the database and Redis
func (r *balanceRepositoryImpl) UpdateBalance(ctx context.Context, userID uint, newBalance float64) error {
	err := r.DB.WithContext(ctx).Model(&model.Balance{}).Where("user_id = ?", userID).Update("balance", newBalance).Error
	if err != nil {
		return err
	}

	return nil
}
