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
		transactionRoutes.POST("/:sessionid", middleware.Authenticate(service.NewJWTService(), "user"), transactionC.MakeTransaction)
		transactionRoutes.GET("", middleware.Authenticate(service.NewJWTService(), "admin"), transactionC.GetAllTransactions)
		transactionRoutes.GET("/:username", middleware.Authenticate(service.NewJWTService(), "admin"), transactionC.GetTransactionsByUsername)
		transactionRoutes.GET("/me", middleware.Authenticate(service.NewJWTService(), "user"), transactionC.GetMyTransactions)
		transactionRoutes.DELETE("/:id", middleware.Authenticate(service.NewJWTService(), "admin"), transactionC.DeleteTransactionByID)
	}
}
