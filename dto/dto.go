package dto

import "walletApp/model"

type TransferRequest struct {
	FromUserID uint    `json:"from_user_id"`
	ToUserID   uint    `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

type TransferResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    map[string]float64 `json:"data"` // debug purpose
}

type DepositRequest struct {
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
}
type DepositResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Balance float64 `json:"balance"`
}

type WithdrawRequest struct {
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
}

type WithdrawResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Balance float64 `json:"balance"`
}

type CheckBalanceRequest struct {
	UserID uint `json:"user_id"`
}

type TransactionHistoryRequest struct {
	UserID uint `json:"user_id"`
}

type TransactionHistoryResponse struct {
	Transactions []model.Transaction `json:"transactions"`
}
