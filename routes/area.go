package routes

import (
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/service"

	"github.com/gin-gonic/gin"
)

func AreaRoutes(router *gin.Engine, areaC controller.AreaController) {
	areaRoutes := router.Group("/api/v1/areas")
	{
		areaRoutes.POST("/", middleware.Authenticate(service.NewJWTService(), "admin"), areaC.CreateArea)
		areaRoutes.GET("/", areaC.GetAllAreas)
		areaRoutes.GET("/:id", middleware.Authenticate(service.NewJWTService(), "admin"), areaC.GetAreaByID)
		areaRoutes.PUT("/:id", middleware.Authenticate(service.NewJWTService(), "admin"), areaC.UpdateAreaByID)
		areaRoutes.DELETE("/:id", middleware.Authenticate(service.NewJWTService(), "admin"), areaC.DeleteAreaByID)
	}
}