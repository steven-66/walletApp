package storage

import (
	"context"
	"walletApp/model"
)

// TransactionRepository defines the interface for transaction-related operations
//
//go:generate mockery --case underscore --name TransactionRepository
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetTransactionsByUserID(ctx context.Context, userID uint) ([]model.Transaction, error)
}
