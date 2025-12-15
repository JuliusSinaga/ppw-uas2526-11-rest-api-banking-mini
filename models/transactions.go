package models

import (
	"time"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AccountID       uint      `json:"account_id"`      // Siapa yang melakukan transaksi
	TargetAccountID *uint     `json:"target_account_id,omitempty"` // Transfer ke siapa (null jika topup)
	TransactionType string    `gorm:"type:varchar(20)" json:"transaction_type"` // "TOPUP", "TRANSFER_IN", "TRANSFER_OUT"
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}