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
		sessionRoutes.DELETE("/:id", middleware.Authenticate(service.NewJWTService(), "admin"), sessionC.DeleteSessionByID)
		
	}

	sessionFilmRoutes := router.Group("/api/v1/sessions/films")
	{
		sessionFilmRoutes.GET("/:filmslug", sessionC.GetSessionsByFilmSlug)
	}

	sessionDetailFilmRoutes := router.Group("/api/v1/sessions/:id/films")
	{
		sessionDetailFilmRoutes.GET("/:filmslug", sessionC.GetSessionDetailByID)
	}
}
