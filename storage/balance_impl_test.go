package storage

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Failed to create mock database: " + err.Error())
	}

	// Wrap the mock database with GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to create GORM DB: " + err.Error())
	}

	return gormDB, mock
}

func TestGetBalance(t *testing.T) {
	db, mock := setupMockDB()
	repo := NewBalanceRepository(db)

	tests := []struct {
		name            string
		userID          uint
		expectedBalance float64
		mockError       error
		expectError     bool
	}{
		{
			name:            "Valid user balance",
			userID:          1,
			expectedBalance: 1000.0,
			mockError:       nil,
			expectError:     false,
		},
		{
			name:            "User not found",
			userID:          2,
			expectedBalance: 0.0,
			mockError:       gorm.ErrRecordNotFound,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the SELECT query
			query := `SELECT .*balance.* FROM "balances" WHERE user_id = \$1 .* LIMIT \$2`
			if tt.mockError == nil {
				mock.ExpectQuery(query).
					WithArgs(tt.userID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(tt.expectedBalance))
			} else {
				mock.ExpectQuery(query).
					WithArgs(tt.userID, 1).
					WillReturnError(tt.mockError)
			}

			ctx := context.Background()
			balance, err := repo.GetBalance(ctx, tt.userID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if balance != tt.expectedBalance {
					t.Errorf("Expected balance %.2f, got %.2f", tt.expectedBalance, balance)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet expectations: %v", err)
			}
		})
	}
}
func TestUpdateBalance(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		newBalance  float64
		mockError   error
		expectError bool
	}{
		{
			name:        "Successful balance update",
			userID:      1,
			newBalance:  150.0,
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "Update fails (user not found)",
			userID:      2,
			newBalance:  200.0,
			mockError:   gorm.ErrRecordNotFound,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock database
			db, mock := setupMockDB()
			repo := NewBalanceRepository(db)

			// Mock transaction Begin
			mock.ExpectBegin()

			query := `UPDATE "balances" SET "balance"=\$1 WHERE user_id = \$2`
			if tt.mockError == nil {
				mock.ExpectExec(query).
					WithArgs(tt.newBalance, tt.userID).
					WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate successful update
			} else {
				mock.ExpectExec(query).
					WithArgs(tt.newBalance, tt.userID).
					WillReturnError(tt.mockError) // Simulate error
			}

			// Mock transaction Commit or Rollback
			if tt.expectError {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			err := repo.UpdateBalance(context.Background(), tt.userID, tt.newBalance)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet expectations: %v", err)
			}
		})
	}
}

func TestNewBalanceRepository(t *testing.T) {
	// Setup mock DB
	gormDB, _ := setupMockDB()

	// Call the function being tested
	repo := NewBalanceRepository(gormDB)

	// Assertions
	assert.NotNil(t, repo)
	assert.IsType(t, &balanceRepositoryImpl{}, repo)
}

func TestNewMockBalanceRepository(t *testing.T) {
	// Call the function being tested with no mock functions
	repo := NewMockBalanceRepository()

	// Assertions
	assert.NotNil(t, repo)
	mockCalled := false
	repo = NewMockBalanceRepository(func(m *mock.Mock) {
		mockCalled = true
		m.On("GetBalance", mock.Anything, uint(1)).Return(100.0, nil)
	})

	// Assertions
	assert.NotNil(t, repo)
	assert.True(t, mockCalled)

	// Test that the mock works as expected
	balance, err := repo.GetBalance(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance)
}
