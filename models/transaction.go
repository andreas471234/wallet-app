package models

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType is string constant for all the available types
type TransactionType string

// TransactionType defination for the TransactionType Field in Transaction
const (
	Debit  TransactionType = "DEBIT"
	Credit TransactionType = "CREDIT"
)

// Transaction model struct that saved in the database
type Transaction struct {
	gorm.Model
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	UserID          uint
	User            User
}

// TransactionDetail detail struct for simplified API response
type TransactionDetail struct {
	ID              int       `json:"id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}

// WalletTransactionRequest struct for sanitize the request in receive and disburse API
type WalletTransactionRequest struct {
	Amount float64 `json:"amount" binding:"required" validate:"gt=0"`
}
