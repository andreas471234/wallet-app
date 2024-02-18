package models

import (
	"fmt"

	"gorm.io/gorm"
)

// User model struct that saved in the database
type User struct {
	gorm.Model
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

// UserDetail detail struct for API response with transaction data
type UserDetail struct {
	User
	Transactions []*TransactionDetail `json:"transactions"`
}

// GetUserDetails compose the user detail from user data
func (u User) GetUserDetails(db *gorm.DB) (UserDetail, error) {
	// Get list of last transactions of this user
	txn_list, err := u.lastTransaction(db)
	if err != nil {
		return UserDetail{}, err
	}

	// Construct the user data combined with last transaction list data
	return UserDetail{
		u,
		txn_list,
	}, nil
}

// lastTransaction getting last 10 transaction for this user
func (u User) lastTransaction(db *gorm.DB) (txns []*TransactionDetail, err error) {
	// Getting the last 10 transaction of the user
	err = db.Model(&Transaction{}).Where("user_id = ?", u.ID).Order("ID DESC").Limit(10).Find(&txns).Error
	return txns, err
}

// ReceiveMoney update the user balance and create transaction
func (u *User) ReceiveMoney(db *gorm.DB, body WalletTransactionRequest) (err error) {
	// Update the user balance
	err = db.Model(&u).Update("balance", u.Balance+body.Amount).Error
	if err != nil {
		return err
	}

	// Create Debit transaction
	err = db.Create(&Transaction{
		Amount:          body.Amount,
		TransactionType: string(Debit),
		UserID:          u.ID,
	}).Error

	return err
}

// DisburseMoney update the user balance and create transaction
func (u *User) DisburseMoney(db *gorm.DB, body WalletTransactionRequest) (err error) {
	// Check if user balance sufficient
	if u.Balance < body.Amount {
		return fmt.Errorf("User Amount Insufficient: %.2f", u.Balance)
	}

	// Update the user balance
	err = db.Model(&u).Update("balance", u.Balance-body.Amount).Error
	if err != nil {
		return err
	}

	// Create Credit transaction
	err = db.Create(&Transaction{
		Amount:          body.Amount,
		TransactionType: string(Credit),
		UserID:          u.ID,
	}).Error

	return err
}
