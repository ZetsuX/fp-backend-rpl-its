package routes

import (
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/service"

	"github.com/gin-gonic/gin"
)

func FilmRoutes(router *gin.Engine, filmC controller.FilmController) {
	filmRoutes := router.Group("/api/v1/movie")
	{
		filmRoutes.POST("/", middleware.Authenticate(service.NewJWTService(), "user"), filmC.CreateFilm)
		filmRoutes.GET("/", filmC.GetAllFilmsAvailable)
		filmRoutes.GET("/all", middleware.Authenticate(service.NewJWTService(), "user"), filmC.GetAllFilms)
		filmRoutes.PUT("/:slug", middleware.Authenticate(service.NewJWTService(), "user"), filmC.UpdateFilm)
		filmRoutes.GET("/:slug", middleware.Authenticate(service.NewJWTService(), "user"), filmC.GetFilmDetailBySlug)
		filmRoutes.DELETE("/:slug", middleware.Authenticate(service.NewJWTService(), "user"), filmC.DeleteFilm)
	}
}
