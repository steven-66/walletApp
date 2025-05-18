package handler

import (
	"context"
	"errors"
	"testing"
	"walletApp/model"
	"walletApp/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTransactionHandler(t *testing.T) {
	handler := NewTransactionHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.TransactionRepo)
}

func TestViewTransactionHistory(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		transactions   []model.Transaction
		repoError      error
		expectSuccess  bool
		expectedLength int
	}{
		{
			name:   "Successful Transaction History Retrieval",
			userID: 1,
			transactions: []model.Transaction{
				{
					ID:     1,
					UserID: 1,
					Amount: 100.0,
					Type:   model.TransactionTypeDeposit,
				},
				{
					ID:     2,
					UserID: 1,
					Amount: -50.0,
					Type:   model.TransactionTypeWithdraw,
				},
			},
			repoError:      nil,
			expectSuccess:  true,
			expectedLength: 2,
		},
		{
			name:           "Empty Transaction History",
			userID:         1,
			transactions:   []model.Transaction{},
			repoError:      nil,
			expectSuccess:  true,
			expectedLength: 0,
		},
		{
			name:           "Repository Error",
			userID:         1,
			transactions:   nil,
			repoError:      errors.New("database error"),
			expectSuccess:  false,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransactionRepo := storage.NewMockTransactionRepository(
				func(m *mock.Mock) {
					m.On("GetTransactionsByUserID", mock.Anything, tt.userID).Return(tt.transactions, tt.repoError)
				},
			)

			handler := &TransactionHandler{
				TransactionRepo: mockTransactionRepo,
			}
			response, err := handler.ViewTransactionHistory(context.Background(), tt.userID)
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedLength, len(response.Transactions))

				// Verify the transactions match
				if tt.expectedLength > 0 {
					for i, tx := range tt.transactions {
						assert.Equal(t, tx.ID, response.Transactions[i].ID)
						assert.Equal(t, tx.UserID, response.Transactions[i].UserID)
						assert.Equal(t, tx.Amount, response.Transactions[i].Amount)
						assert.Equal(t, tx.Type, response.Transactions[i].Type)
					}
				}
			} else {
				assert.Error(t, err)
				assert.Nil(t, response)
			}
		})
	}
}
