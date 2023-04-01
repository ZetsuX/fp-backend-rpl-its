package routes

import (
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/service"

	"github.com/gin-gonic/gin"
)

func FilmRoutes(router *gin.Engine, filmC controller.FilmController) {
	filmRoutes := router.Group("/api/v1/films")
	{
		filmRoutes.POST("", middleware.Authenticate(service.NewJWTService(), "admin"), filmC.CreateFilm)
		filmRoutes.GET("", filmC.GetAllFilmsNowPlaying)
		filmRoutes.GET("/coming-soon", filmC.GetAllFilmsComingSoon)
		filmRoutes.GET("/all", filmC.GetAllFilms)
		filmRoutes.PUT("/:slug", middleware.Authenticate(service.NewJWTService(), "admin"), filmC.UpdateFilm)
		filmRoutes.GET("/:slug", filmC.GetFilmDetailBySlug)
		filmRoutes.DELETE("/:slug", middleware.Authenticate(service.NewJWTService(), "admin"), filmC.DeleteFilm)
	}
}
