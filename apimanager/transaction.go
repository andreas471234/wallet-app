package apimanager

import (
	"math"
	"net/http"
	"wallet-service-gin/models"
	"wallet-service-gin/utils"

	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"gorm.io/gorm"
)

// Galidator are used to custom the validation message for API Requests for better understanding
var (
	g = galidator.New().CustomMessages(galidator.Messages{
		"required": "$field is required",
	})
	customizer = g.Validator(models.WalletTransactionRequest{})
)

// GetUsers responds with the list of all users as JSON and apply the filter also
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get all the query params for filter and pagination
		name := c.Query("name")
		min_balance := c.Query("min_balance")
		max_balance := c.Query("max_balance")
		page := c.Query("page")
		page_size := c.Query("page_size")

		// Sanitize the page_size and page of the pagination
		start, count, curr_page := utils.GetLimitOffset(page_size, 10, page, 1)

		// Create the custom filter and sanitize the query param data
		var rawFilters [][]string
		rawFilters = models.AddRawFilter("name", name, rawFilters)
		rawFilters = models.AddRawFilter("min_balance", min_balance, rawFilters)
		rawFilters = models.AddRawFilter("max_balance", max_balance, rawFilters)

		// Find all the user with applied filter and pagination
		res, err := models.FindUserBy(db, []models.User{}, rawFilters, start, count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check for total count of users without filters
		var total_count int64
		db.Model(models.User{}).Count(&total_count)

		// Create the response data with metadata for pagination purposes
		response := gin.H{
			"data": res,
			"meta_data": map[string]interface{}{
				"page":          curr_page,
				"page_size":     count,
				"total_records": total_count,
				"total_page":    int(math.Ceil(float64(total_count) / float64(count))),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetUserDetails locates the user whose ID value matches the id
// parameter sent by the client, then returns that user details as a response.
func GetUserDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find User with given ID from the param or return error if not found
		user, err := models.FindUserID(db, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Construct the user details data with transactions list for API response
		details, err := user.GetUserDetails(db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, details)
	}
}

// Receive collect money from user account
func Receive(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find User with given ID from the param or return error if not found
		user, err := models.FindUserID(db, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// using BindJson method to serialize body with struct and give validation error
		body := models.WalletTransactionRequest{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": customizer.DecryptErrors(err)})
			return
		}

		// Create db transaction so we can rollback in case something happen in the middle of update
		tx := db.Begin()

		// Update the user balance and create Debit transactions
		err = user.ReceiveMoney(tx, body)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Committing the all db transactions changes into db and return error if failed
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Contruct the user detail data to show the new balance and transactions
		details, err := user.GetUserDetails(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, details)
	}
}

// Disburse disburse money from user account
func Disburse(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find User with given ID from the param or return error if not found
		user, err := models.FindUserID(db, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// using BindJson method to serialize body with struct and give validation error
		body := models.WalletTransactionRequest{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": customizer.DecryptErrors(err)})
			return
		}

		// Create db transaction so we can rollback in case something happen in the middle of update
		tx := db.Begin()

		// Update the user balance and create Credit transactions
		err = user.DisburseMoney(tx, body)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Committing the all db transactions changes into db and return error if failed
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Contruct the user detail data to show the new balance and transactions
		details, err := user.GetUserDetails(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, details)
	}
}

// GetTransactionList locates the user whose ID value matches the id
// parameter sent by the client, then returns list of transactions of the user with applied filter
func GetTransactionList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find User with given ID from the param or return error if not found
		user, err := models.FindUserID(db, c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get all the query params for filter and pagination
		txn_type := c.Query("type")
		page := c.Query("page")
		page_size := c.Query("page_size")

		// Sanitize the page_size and page of the pagination
		start, count, curr_page := utils.GetLimitOffset(page_size, 10, page, 1)

		// Create the transaction filter for specific user
		filter := []models.Transaction{{UserID: user.ID}}

		// Create the custom filter and sanitize the query param data
		var rawFilters [][]string
		rawFilters = models.AddRawFilter("type", txn_type, rawFilters)

		// Find all the transaction of particular user with applied filter and pagination
		res, err := models.FindTransactionBy(db, filter, rawFilters, start, count)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check for total count of transactions of the users without filters
		var total_count int64
		db.Model(models.Transaction{}).Where("user_id = ?", user.ID).Count(&total_count)

		// Create the response data with metadata for pagination purposes
		response := gin.H{
			"data": res,
			"meta_data": map[string]interface{}{
				"page":          curr_page,
				"page_size":     count,
				"total_records": total_count,
				"total_page":    int(math.Ceil(float64(total_count) / float64(count))),
			},
		}

		c.JSON(http.StatusOK, response)
	}
}
