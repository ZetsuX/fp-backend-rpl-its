package main

import (
	"fmt"
	"fp-rpl/config"
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/repository"
	"fp-rpl/routes"
	"fp-rpl/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	// Setting Up Database
	db := config.DBSetup()

	// Setting Up Repositories
	userR := repository.NewUserRepository(db)
	filmR := repository.NewFilmRepository(db)
	areaR := repository.NewAreaRepository(db)
	sessionR := repository.NewSessionRepository(db)
	spotR := repository.NewSpotRepository(db)

	// Setting Up Services
	userS := service.NewUserService(userR)
	filmS := service.NewFilmService(filmR)
	jwtS := service.NewJWTService()
	areaS := service.NewAreaService(areaR)
	sessionS := service.NewSessionService(sessionR, spotR)

	// Setting Up Controllers
	userC := controller.NewUserController(userS, jwtS)
	filmC := controller.NewFilmController(filmS)
	areaC := controller.NewAreaController(areaS)
	sessionC := controller.NewSessionController(sessionS, areaS)

	defer config.DBClose(db)

	// Setting Up Server
	server := gin.Default()
	server.Use(
		middleware.CORSMiddleware(),
	)

	// Setting Up Routes
	routes.UserRoutes(server, userC)
	routes.FilmRoutes(server,filmC)
	routes.AreaRoutes(server, areaC)
	routes.SessionRoutes(server, sessionC)

	// Running in localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
