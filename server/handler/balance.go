package handler

import (
	"context"
	"fmt"
	"log"
	"walletApp/config"
	"walletApp/dto"
	"walletApp/model"
	"walletApp/storage"
)

type BalanceHandler struct {
	BalanceRepo     storage.BalanceRepository
	TransactionRepo storage.TransactionRepository
}

// NewBalanceHandler creates a new instance of BalanceHandler
func NewBalanceHandler() *BalanceHandler {
	return &BalanceHandler{BalanceRepo: storage.NewBalanceRepository(config.DB), TransactionRepo: storage.NewTransactionRepository(config.DB)}
}

func (c *BalanceHandler) CheckBalance(ctx context.Context, userID uint) (float64, error) {
	balance, err := c.BalanceRepo.GetBalance(ctx, userID)
	if err != nil {
		log.Printf("Error fetching balance for user %d: %v\n", userID, err)
		return 0, fmt.Errorf("failed to fetch balance for user %d: %w", userID, err)
	}

	return balance, nil
}

func (c *BalanceHandler) Deposit(ctx context.Context, request *dto.DepositRequest) (*dto.DepositResponse, error) {
	userID, amount := request.UserID, request.Amount
	// Fetch balance
	balance, err := c.BalanceRepo.GetBalance(ctx, userID)
	if err != nil {
		log.Printf("Error fetching balance for user %d: %v\n", userID, err)
		return nil, fmt.Errorf("failed to fetch balance for user %d: %w", userID, err)
	}

	// Update balance
	newBalance := balance + amount
	err = c.BalanceRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		log.Printf("Error updating balance for user %d: %v\n", userID, err)
		return nil, fmt.Errorf("failed to update balance for user %d: %w", userID, err)
	}

	// Create deposit transaction
	transaction := &model.Transaction{
		UserID: userID,
		Amount: amount,
		Type:   model.TransactionTypeDeposit,
	}
	err = c.TransactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Printf("Error creating transaction for user %d: %v\n", userID, err)
		return nil, fmt.Errorf("failed to create transaction for user %d: %w", userID, err)
	}

	return &dto.DepositResponse{
		Success: true,
		Message: "Success Deposit",
		Balance: newBalance,
	}, nil
}

func (c *BalanceHandler) Withdraw(ctx context.Context, request *dto.WithdrawRequest) (*dto.WithdrawResponse, error) {
	userID, amount := request.UserID, request.Amount
	// Fetch balance
	balance, err := c.BalanceRepo.GetBalance(ctx, request.UserID)
	if err != nil {
		log.Printf("Error fetching balance for user %d: %v\n", userID, err)
		return nil, fmt.Errorf("failed to fetch balance for user %d: %w", userID, err)
	}

	// Check if balance is sufficient
	if balance < amount {
		log.Printf("Insufficient balance for user %d\n", userID)
		return nil, fmt.Errorf("insufficient balance for user %d", userID)
	}

	// Update balance
	newBalance := balance - amount
	err = c.BalanceRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		log.Printf("Error updating balance for user %d: %v\n", userID, err)
		return nil, fmt.Errorf("failed to update balance for user %d: %w", userID, err)
	}

	// Create withdraw transaction
	transaction := &model.Transaction{
		UserID: userID,
		Amount: amount,
		Type:   model.TransactionTypeWithdraw,
	}
	err = c.TransactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Printf("Error creating transaction for user %d: %v\n", userID, err)
		return nil, fmt.Errorf("failed to create transaction for user %d: %w", userID, err)
	}

	return &dto.WithdrawResponse{
		Success: true,
		Message: "Withdrawal successful",
		Balance: newBalance,
	}, nil
}

func (c *BalanceHandler) Transfer(ctx context.Context, request *dto.TransferRequest) (*dto.TransferResponse, error) {
	fromUserID, toUserID, amount := request.FromUserID, request.ToUserID, request.Amount
	// Fetch sender's balance
	senderBalance, err := c.BalanceRepo.GetBalance(ctx, fromUserID)
	if err != nil {
		log.Printf("Error fetching balance for sender %d: %v\n", fromUserID, err)
		return &dto.TransferResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to fetch balance for sender %d", fromUserID),
		}, err
	}

	if senderBalance < amount {
		log.Printf("Insufficient balance for sender %d\n", fromUserID)
		return &dto.TransferResponse{
			Success: false,
			Message: "Insufficient balance",
		}, fmt.Errorf("insufficient balance for sender %d", fromUserID)
	}

	// Fetch recipient's balance
	recipientBalance, err := c.BalanceRepo.GetBalance(ctx, toUserID)
	if err != nil {
		log.Printf("Error fetching balance for recipient %d: %v\n", toUserID, err)
		return &dto.TransferResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to fetch balance for recipient %d", toUserID),
		}, err
	}

	// Update balances
	newSenderBalance := senderBalance - amount
	newRecipientBalance := recipientBalance + amount

	err = c.BalanceRepo.UpdateBalance(ctx, fromUserID, newSenderBalance)
	if err != nil {
		log.Printf("Error updating balance for sender %d: %v\n", fromUserID, err)
		return &dto.TransferResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update balance for sender %d", fromUserID),
		}, err
	}

	err = c.BalanceRepo.UpdateBalance(ctx, toUserID, newRecipientBalance)
	if err != nil {
		log.Printf("Error updating balance for recipient %d: %v\n", toUserID, err)
		return &dto.TransferResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to update balance for recipient %d", toUserID),
		}, err
	}

	// Log transactions for both sender and recipient
	senderTransaction := model.Transaction{
		UserID: fromUserID,
		Type:   model.TransactionTypeTransferSend,
		Amount: -amount,
	}
	err = c.TransactionRepo.CreateTransaction(ctx, &senderTransaction)
	if err != nil {
		log.Printf("Error logging transaction for sender %d: %v\n", fromUserID, err)
	}

	recipientTransaction := model.Transaction{
		UserID: toUserID,
		Type:   model.TransactionTypeTransferReceive,
		Amount: amount,
	}
	err = c.TransactionRepo.CreateTransaction(ctx, &recipientTransaction)
	if err != nil {
		log.Printf("Error logging transaction for recipient %d: %v\n", toUserID, err)
	}

	return &dto.TransferResponse{
		Success: true,
		Message: "Transfer successful",
		Data: map[string]float64{
			"sender_balance":    newSenderBalance,
			"recipient_balance": newRecipientBalance,
		},
	}, nil
}
