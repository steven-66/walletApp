package storage

import (
	"context"
)

// BalanceRepository defines the interface for balance-related operations
//
//go:generate mockery  --case underscore --name BalanceRepository
type BalanceRepository interface {
	GetBalance(ctx context.Context, userID uint) (float64, error)
	UpdateBalance(ctx context.Context, userID uint, newBalance float64) error
}
