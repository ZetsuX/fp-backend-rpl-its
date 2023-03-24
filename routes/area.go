package routes

import (
	"fp-rpl/controller"

	"github.com/gin-gonic/gin"
)

func AreaRoutes(router *gin.Engine, areaC controller.AreaController) {
	areaRoutes := router.Group("/api/v1/area")
	{
		areaRoutes.POST("/", areaC.CreateArea)
		areaRoutes.GET("/", areaC.GetAllAreas)
	}
}