package model

import "time"

type Transaction struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `json:"user_id"`
	Type      TransactionType `json:"type"` // Deposit, Withdraw, Transfer
	Amount    float64         `json:"amount"`
	Timestamp time.Time       `gorm:"autoCreateTime" json:"timestamp"`
}

type TransactionType uint16

const (
	TransactionTypeDeposit TransactionType = iota
	TransactionTypeWithdraw
	TransactionTypeTransferSend
	TransactionTypeTransferReceive
)

func (t TransactionType) String() string {
	switch t {
	case TransactionTypeDeposit:
		return "Deposit"
	case TransactionTypeWithdraw:
		return "Withdraw"
	case TransactionTypeTransferSend:
		return "TransferSend"
	case TransactionTypeTransferReceive:
		return "TransferReceive"
	default:
		return "Unknown"
	}
}
