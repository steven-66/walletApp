package storage

import (
	"context"
	"errors"
	"testing"
	"time"
	"walletApp/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionds(t *testing.T) {
	fixedTime := time.Now()

	tests := []struct {
		name        string
		transaction *model.Transaction
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful Transaction Creation",
			transaction: &model.Transaction{
				UserID:    1,
				Amount:    100.0,
				Type:      model.TransactionTypeDeposit,
				Timestamp: fixedTime,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "transactions"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Database Error",
			transaction: &model.Transaction{
				UserID:    2,
				Amount:    -50.0,
				Type:      model.TransactionTypeWithdraw,
				Timestamp: fixedTime,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "transactions"`).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB
			gormDB, mock := setupMockDB()
			tt.setupMock(mock)

			repo := &TransactionRepositoryImpl{DB: gormDB}
			err := repo.CreateTransaction(context.Background(), tt.transaction)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet expectations: %v", err)
			}
		})
	}
}

func TestGetTransactionsByUserID(t *testing.T) {
	// Define a fixed timestamp for testing
	fixedTime := time.Now()

	tests := []struct {
		name            string
		userID          uint
		setupMock       func(sqlmock.Sqlmock)
		expectedResults []model.Transaction
		expectError     bool
	}{
		{
			name:   "Successful Transactions Retrieval",
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "type", "timestamp"}).
					AddRow(1, 1, 100.0, model.TransactionTypeDeposit, fixedTime).
					AddRow(2, 1, -50.0, model.TransactionTypeWithdraw, fixedTime)
				mock.ExpectQuery(`SELECT.*FROM "transactions".*WHERE.*user_id = \$1.*ORDER BY timestamp DESC`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedResults: []model.Transaction{
				{
					ID:        1,
					UserID:    1,
					Amount:    100.0,
					Type:      model.TransactionTypeDeposit,
					Timestamp: fixedTime,
				},
				{
					ID:        2,
					UserID:    1,
					Amount:    -50.0,
					Type:      model.TransactionTypeWithdraw,
					Timestamp: fixedTime,
				},
			},
			expectError: false,
		},
		{
			name:   "Empty Transactions",
			userID: 2,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "type", "timestamp"})
				mock.ExpectQuery(`SELECT.*FROM "transactions".*WHERE.*user_id = \$1.*ORDER BY timestamp DESC`).
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedResults: []model.Transaction{},
			expectError:     false,
		},
		{
			name:   "Database Error",
			userID: 3,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT.*FROM "transactions".*WHERE.*user_id = \$1.*ORDER BY timestamp DESC`).
					WithArgs(3).
					WillReturnError(errors.New("database connection error"))
			},
			expectedResults: nil,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB
			gormDB, mock := setupMockDB()
			tt.setupMock(mock)

			// Create repository with mocked DB
			repo := &TransactionRepositoryImpl{DB: gormDB}

			// Call the method being tested
			transactions, err := repo.GetTransactionsByUserID(context.Background(), tt.userID)

			// Assertions
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, transactions)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResults), len(transactions))

				// Verify the transactions match
				if len(tt.expectedResults) > 0 {
					for i, tx := range tt.expectedResults {
						assert.Equal(t, tx.ID, transactions[i].ID)
						assert.Equal(t, tx.UserID, transactions[i].UserID)
						assert.Equal(t, tx.Amount, transactions[i].Amount)
						assert.Equal(t, tx.Type, transactions[i].Type)
						assert.Equal(t, tx.Timestamp.Unix(), transactions[i].Timestamp.Unix())
					}
				}
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNewTransactionRepository(t *testing.T) {
	gormDB, _ := setupMockDB()
	repo := NewTransactionRepository(gormDB)

	assert.NotNil(t, repo)
	assert.IsType(t, &TransactionRepositoryImpl{}, repo)
}

func TestNewMockTransactionRepository(t *testing.T) {
	repo := NewMockTransactionRepository()
	assert.NotNil(t, repo)

	mockCalled := false
	repo = NewMockTransactionRepository(func(m *mock.Mock) {
		mockCalled = true
		m.On("GetTransactionsByUserID", mock.Anything, uint(1)).Return([]model.Transaction{}, nil)
	})

	assert.NotNil(t, repo)
	assert.True(t, mockCalled)

	// Test that the mock works as expected
	transactions, err := repo.GetTransactionsByUserID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Empty(t, transactions)
}
