package models

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// FindUserID finding user from given params id
func FindUserID(db *gorm.DB, id string) (user *User, err error) {
	// Sanitize the id input param should be integer
	if _, err := strconv.Atoi(id); err != nil {
		return nil, fmt.Errorf("id given is not a string")
	}

	// Get the user with this id
	err = db.First(&user, id).Error
	return user, err
}

// FindUserBy gives list of user from applied filter with offset and limit
func FindUserBy(db *gorm.DB, finders []User, rawFinders [][]string, offset, limit int) (users []*User, err error) {
	// Applied the model specific filters
	for _, f := range finders {
		db = db.Where(f)
	}

	// Applied the custom filters
	for _, f := range rawFinders {
		db = db.Where(f[0], f[1:])
	}

	// Get the data with offset and limit
	err = db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// FindTransactionBy gives list of transaction from applied filter with offset and limit
func FindTransactionBy(db *gorm.DB, finders []Transaction, rawFinders [][]string, offset, limit int) (txns []*TransactionDetail, err error) {
	// Applied the model specific filters
	for _, f := range finders {
		db = db.Where(f)
	}

	// Applied the custom filters
	for _, f := range rawFinders {
		db = db.Where(f[0], f[1:])
	}

	// Get the data with offset and limit
	err = db.Model(Transaction{}).Offset(offset).Limit(limit).Find(&txns).Error
	return txns, err
}

// AddRawFilter sanitize filter and build custom raw filter mapping with columns for user model
func AddRawFilter(field string, value string, filters [][]string) [][]string {
	// Check if value not empty
	if len(value) > 0 {
		filter := []string{}
		switch field {
		case "name":
			filter = []string{"name like ?", "%" + value + "%"}

		case "type":
			// Check only if value match with CREDIT and DEBIT, else not applying this filter
			if value == string(Credit) {
				filter = []string{"transaction_type = ?", value}
			} else if value == string(Debit) {
				filter = []string{"transaction_type = ?", value}
			} else {
				return filters
			}

		case "min_balance":
			// Check if given value not integer then skipping this filter
			if _, err := strconv.Atoi(value); err != nil {
				return filters
			}
			filter = []string{"balance >= ?", value}

		case "max_balance":
			// Check if given value not integer then skipping this filter
			if _, err := strconv.Atoi(value); err != nil {
				return filters
			}
			filter = []string{"balance <= ?", value}
		}

		// Append the custom filter to list of filters after sanitizing it
		filters = append(filters, filter)
	}
	return filters
}
