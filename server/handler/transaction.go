package handler

import (
	"context"
	"fmt"
	"log"
	"walletApp/config"
	"walletApp/dto"
	"walletApp/storage"
)

type TransactionHandler struct {
	TransactionRepo storage.TransactionRepository
}

// NewTransactionHandler creates a new instance of TransactionHandler
func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{TransactionRepo: storage.NewTransactionRepository(config.DB)}
}

func (c *TransactionHandler) ViewTransactionHistory(ctx context.Context, userID uint) (*dto.TransactionHistoryResponse, error) {
	transactions, err := c.TransactionRepo.GetTransactionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error fetching transaction history for user %d: %v\n", userID, err)
		return nil, fmt.Errorf("failed to fetch transaction history for user %d: %w", userID, err)
	}

	return &dto.TransactionHistoryResponse{
		Transactions: transactions,
	}, nil
}
