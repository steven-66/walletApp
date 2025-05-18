package handler

import (
	"context"
	"errors"
	"testing"
	"walletApp/dto"
	"walletApp/model"
	"walletApp/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBalanceHandler(t *testing.T) {
	handler := NewBalanceHandler()
	assert.NotNil(t, handler, "Expected non-nil handler, got nil")
	assert.NotNil(t, handler.BalanceRepo, "Expected non-nil BalanceRepo, got nil")
	assert.NotNil(t, handler.TransactionRepo, "Expected non-nil TransactionRepo, got nil")
}

func TestCheckBalance(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		mockBalance   float64
		mockError     error
		expectedValue float64
		expectError   bool
	}{
		{
			name:          "Success",
			userID:        1,
			mockBalance:   100.0,
			mockError:     nil,
			expectedValue: 100.0,
			expectError:   false,
		},
		{
			name:          "Repository Error",
			userID:        1,
			mockBalance:   0.0,
			mockError:     errors.New("database error"),
			expectedValue: 0.0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceRepo := storage.NewMockBalanceRepository(func(mocker *mock.Mock) {
				mocker.On("GetBalance", mock.Anything, tt.userID).Return(tt.mockBalance, tt.mockError)
			})

			handler := &BalanceHandler{
				BalanceRepo: mockBalanceRepo,
			}

			balance, err := handler.CheckBalance(context.Background(), tt.userID)

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if balance != tt.expectedValue {
				t.Errorf("Expected balance %f, got %f", tt.expectedValue, balance)
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		name               string
		request            *dto.DepositRequest
		initialBalance     float64
		getBalanceError    error
		updateBalanceError error
		createTxError      error
		expectSuccess      bool
		expectedBalance    float64
	}{
		{
			name: "Successful Deposit",
			request: &dto.DepositRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     100.0,
			getBalanceError:    nil,
			updateBalanceError: nil,
			createTxError:      nil,
			expectSuccess:      true,
			expectedBalance:    150.0,
		},
		{
			name: "GetBalance Error",
			request: &dto.DepositRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     0.0,
			getBalanceError:    errors.New("database error"),
			updateBalanceError: nil,
			createTxError:      nil,
			expectSuccess:      false,
			expectedBalance:    0.0,
		},
		{
			name: "UpdateBalance Error",
			request: &dto.DepositRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     100.0,
			getBalanceError:    nil,
			updateBalanceError: errors.New("update error"),
			createTxError:      nil,
			expectSuccess:      false,
			expectedBalance:    0.0,
		},
		{
			name: "CreateTransaction Error",
			request: &dto.DepositRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     100.0,
			getBalanceError:    nil,
			updateBalanceError: nil,
			createTxError:      errors.New("transaction error"),
			expectSuccess:      false,
			expectedBalance:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceRepo := storage.NewMockBalanceRepository(func(mocker *mock.Mock) {
				mocker.On("GetBalance", mock.Anything, tt.request.UserID).Return(tt.initialBalance, tt.getBalanceError)
				mocker.On("UpdateBalance", mock.Anything, tt.request.UserID, tt.request.Amount+tt.initialBalance).Return(tt.updateBalanceError)
			})

			mockTxRepo := storage.NewMockTransactionRepository(func(mocker *mock.Mock) {
				mocker.On("CreateTransaction", mock.Anything, mock.Anything).Return(tt.createTxError)
			})

			handler := &BalanceHandler{
				BalanceRepo:     mockBalanceRepo,
				TransactionRepo: mockTxRepo,
			}

			response, err := handler.Deposit(context.Background(), tt.request)

			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if response == nil {
					t.Error("Expected non-nil response, got nil")
				} else {
					if !response.Success {
						t.Error("Expected success=true, got false")
					}
					if response.Balance != tt.expectedBalance {
						t.Errorf("Expected balance %f, got %f", tt.expectedBalance, response.Balance)
					}
				}
			} else {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name               string
		request            *dto.WithdrawRequest
		initialBalance     float64
		getBalanceError    error
		updateBalanceError error
		createTxError      error
		expectSuccess      bool
		expectedBalance    float64
	}{
		{
			name: "Successful Withdrawal",
			request: &dto.WithdrawRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     100.0,
			getBalanceError:    nil,
			updateBalanceError: nil,
			createTxError:      nil,
			expectSuccess:      true,
			expectedBalance:    50.0,
		},
		{
			name: "Insufficient Balance",
			request: &dto.WithdrawRequest{
				UserID: 1,
				Amount: 150.0,
			},
			initialBalance:     100.0,
			getBalanceError:    nil,
			updateBalanceError: nil,
			createTxError:      nil,
			expectSuccess:      false,
			expectedBalance:    0.0,
		},
		{
			name: "GetBalance Error",
			request: &dto.WithdrawRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     0.0,
			getBalanceError:    errors.New("database error"),
			updateBalanceError: nil,
			createTxError:      nil,
			expectSuccess:      false,
			expectedBalance:    0.0,
		},
		{
			name: "UpdateBalance Error",
			request: &dto.WithdrawRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     100.0,
			getBalanceError:    nil,
			updateBalanceError: errors.New("update error"),
			createTxError:      nil,
			expectSuccess:      false,
			expectedBalance:    0.0,
		},
		{
			name: "CreateTransaction Error",
			request: &dto.WithdrawRequest{
				UserID: 1,
				Amount: 50.0,
			},
			initialBalance:     100.0,
			getBalanceError:    nil,
			updateBalanceError: nil,
			createTxError:      errors.New("transaction error"),
			expectSuccess:      false,
			expectedBalance:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceRepo := storage.NewMockBalanceRepository(func(mocker *mock.Mock) {
				mocker.On("GetBalance", mock.Anything, tt.request.UserID).Return(tt.initialBalance, tt.getBalanceError)
				mocker.On("UpdateBalance", mock.Anything, tt.request.UserID, tt.initialBalance-tt.request.Amount).Return(tt.updateBalanceError)
			})

			mockTxRepo := storage.NewMockTransactionRepository(func(mocker *mock.Mock) {
				mocker.On("CreateTransaction", mock.Anything, mock.Anything).Return(tt.createTxError)
			})

			handler := &BalanceHandler{
				BalanceRepo:     mockBalanceRepo,
				TransactionRepo: mockTxRepo,
			}

			response, err := handler.Withdraw(context.Background(), tt.request)

			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if response == nil {
					t.Error("Expected non-nil response, got nil")
				} else {
					if !response.Success {
						t.Error("Expected success=true, got false")
					}
					if response.Balance != tt.expectedBalance {
						t.Errorf("Expected balance %f, got %f", tt.expectedBalance, response.Balance)
					}
				}
			} else {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	tests := []struct {
		name                     string
		request                  *dto.TransferRequest
		senderBalance            float64
		recipientBalance         float64
		getSenderError           error
		getRecipientError        error
		updateSenderError        error
		updateRecipientError     error
		createSenderTxError      error
		createRecipientTxError   error
		expectSuccess            bool
		expectedSenderBalance    float64
		expectedRecipientBalance float64
	}{
		{
			name: "Successful Transfer",
			request: &dto.TransferRequest{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     50.0,
			},
			senderBalance:            100.0,
			recipientBalance:         100.0,
			getSenderError:           nil,
			getRecipientError:        nil,
			updateSenderError:        nil,
			updateRecipientError:     nil,
			createSenderTxError:      nil,
			createRecipientTxError:   nil,
			expectSuccess:            true,
			expectedSenderBalance:    50.0,
			expectedRecipientBalance: 150.0,
		},
		{
			name: "Insufficient Sender Balance",
			request: &dto.TransferRequest{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     150.0,
			},
			senderBalance:            100.0,
			recipientBalance:         100.0,
			getSenderError:           nil,
			getRecipientError:        nil,
			updateSenderError:        nil,
			updateRecipientError:     nil,
			createSenderTxError:      nil,
			createRecipientTxError:   nil,
			expectSuccess:            false,
			expectedSenderBalance:    100.0,
			expectedRecipientBalance: 100.0,
		},
		{
			name: "Get Sender Balance Error",
			request: &dto.TransferRequest{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     50.0,
			},
			senderBalance:            0.0,
			recipientBalance:         100.0,
			getSenderError:           errors.New("database error"),
			getRecipientError:        nil,
			updateSenderError:        nil,
			updateRecipientError:     nil,
			createSenderTxError:      nil,
			createRecipientTxError:   nil,
			expectSuccess:            false,
			expectedSenderBalance:    0.0,
			expectedRecipientBalance: 100.0,
		},
		{
			name: "Get Recipient Balance Error",
			request: &dto.TransferRequest{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     50.0,
			},
			senderBalance:            100.0,
			recipientBalance:         0.0,
			getSenderError:           nil,
			getRecipientError:        errors.New("database error"),
			updateSenderError:        nil,
			updateRecipientError:     nil,
			createSenderTxError:      nil,
			createRecipientTxError:   nil,
			expectSuccess:            false,
			expectedSenderBalance:    100.0,
			expectedRecipientBalance: 0.0,
		},
		{
			name: "Update Sender Balance Error",
			request: &dto.TransferRequest{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     50.0,
			},
			senderBalance:            100.0,
			recipientBalance:         100.0,
			getSenderError:           nil,
			getRecipientError:        nil,
			updateSenderError:        errors.New("update error"),
			updateRecipientError:     nil,
			createSenderTxError:      nil,
			createRecipientTxError:   nil,
			expectSuccess:            false,
			expectedSenderBalance:    100.0,
			expectedRecipientBalance: 100.0,
		},
		{
			name: "Update Recipient Balance Error",
			request: &dto.TransferRequest{
				FromUserID: 1,
				ToUserID:   2,
				Amount:     50.0,
			},
			senderBalance:            100.0,
			recipientBalance:         100.0,
			getSenderError:           nil,
			getRecipientError:        nil,
			updateSenderError:        nil,
			updateRecipientError:     errors.New("update error"),
			createSenderTxError:      nil,
			createRecipientTxError:   nil,
			expectSuccess:            false,
			expectedSenderBalance:    100.0,
			expectedRecipientBalance: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBalanceRepo := storage.NewMockBalanceRepository(
				func(m *mock.Mock) {
					// Set up expectations for sender balance check
					m.On("GetBalance", mock.Anything, tt.request.FromUserID).Return(tt.senderBalance, tt.getSenderError)
					m.On("GetBalance", mock.Anything, tt.request.ToUserID).Return(tt.recipientBalance, tt.getRecipientError)
					m.On("UpdateBalance", mock.Anything, tt.request.FromUserID, tt.senderBalance-tt.request.Amount).Return(tt.updateSenderError)
					m.On("UpdateBalance", mock.Anything, tt.request.ToUserID, tt.recipientBalance+tt.request.Amount).Return(tt.updateRecipientError)

				},
			)

			mockTransactionRepo := storage.NewMockTransactionRepository(
				func(m *mock.Mock) {
					if tt.getSenderError == nil && tt.getRecipientError == nil &&
						tt.updateSenderError == nil && tt.updateRecipientError == nil &&
						tt.senderBalance >= tt.request.Amount {

						m.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *model.Transaction) bool {
							return tx.UserID == tt.request.FromUserID &&
								tx.Type == model.TransactionTypeTransferSend &&
								tx.Amount == -tt.request.Amount
						})).Return(tt.createSenderTxError)

						m.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *model.Transaction) bool {
							return tx.UserID == tt.request.ToUserID &&
								tx.Type == model.TransactionTypeTransferReceive &&
								tx.Amount == tt.request.Amount
						})).Return(tt.createRecipientTxError)
					}
				},
			)

			handler := &BalanceHandler{
				BalanceRepo:     mockBalanceRepo,
				TransactionRepo: mockTransactionRepo,
			}

			response, err := handler.Transfer(context.Background(), tt.request)

			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				assert.Equal(t, "Transfer successful", response.Message)

				if response.Data != nil {
					assert.Equal(t, tt.expectedSenderBalance, response.Data["sender_balance"])
					assert.Equal(t, tt.expectedRecipientBalance, response.Data["recipient_balance"])
				}
			} else {
				assert.Error(t, err)
				if response != nil {
					assert.False(t, response.Success)
				}
			}
		})
	}
}
