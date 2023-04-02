package routes

import (
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/service"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(router *gin.Engine, transactionC controller.TransactionController) {
	transactionRoutes := router.Group("/api/v1/transactions")
	{
		transactionRoutes.GET("", middleware.Authenticate(service.NewJWTService(), "admin"), transactionC.GetAllTransactions)
		transactionRoutes.GET("/me", middleware.Authenticate(service.NewJWTService(), "user"), transactionC.GetMyTransactions)
		transactionRoutes.DELETE("/:id", middleware.Authenticate(service.NewJWTService(), "admin"), transactionC.DeleteTransactionByID)
	}

	transactionUserRoutes := router.Group("/api/v1/transactions/users")
	{
		transactionUserRoutes.GET("/:username", middleware.Authenticate(service.NewJWTService(), "admin"), transactionC.GetTransactionsByUsername)
	}

	transactionSessionRoutes := router.Group("/api/v1/transactions/sessions")
	{
		transactionSessionRoutes.POST("/:sessionid", middleware.Authenticate(service.NewJWTService(), "user"), transactionC.MakeTransaction)
	}
}
