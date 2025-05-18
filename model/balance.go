package model

import "time"

type Balance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex" json:"user_id"` // Ensures one balance per user
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
