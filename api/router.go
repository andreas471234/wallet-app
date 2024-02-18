package api

import (
	"wallet-service-gin/apimanager"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// API all the router that exposed from the server
func API(db *gorm.DB) {
	router := gin.Default()

	// Router Group for users
	users := router.Group("/users")
	{
		users.GET("/", apimanager.GetUsers(db))
		users.GET("/:id", apimanager.GetUserDetails(db))
	}

	// Router Group for wallets
	wallets := router.Group("/wallets")
	{
		wallets.GET("/:id", apimanager.GetTransactionList(db))
		wallets.POST("/receive/:id", apimanager.Receive(db))
		wallets.POST("/disburse/:id", apimanager.Disburse(db))
	}

	router.Run("localhost:8080")
}
