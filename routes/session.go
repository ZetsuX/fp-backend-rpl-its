package routes

import (
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/service"

	"github.com/gin-gonic/gin"
)

func SessionRoutes(router *gin.Engine, sessionC controller.SessionController) {
	sessionRoutes := router.Group("/api/v1/sessions")
	{
		sessionRoutes.POST("", middleware.Authenticate(service.NewJWTService(), "admin"), sessionC.CreateSession)
		sessionRoutes.GET("", middleware.Authenticate(service.NewJWTService(), "admin"), sessionC.GetAllSessions)
		sessionRoutes.GET("/:filmslug", sessionC.GetSessionsByFilmSlug)
		sessionRoutes.DELETE("/:id", middleware.Authenticate(service.NewJWTService(), "admin"), sessionC.DeleteSessionByID)
		sessionRoutes.GET("/:filmslug/:id", sessionC.GetSessionDetailByID)
	}
}
