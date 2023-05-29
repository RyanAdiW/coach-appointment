package main

import (
	"fita/project/coach-appointment/config"
	"fita/project/coach-appointment/database"
	appointmentController "fita/project/coach-appointment/delivery/controller/appointment"
	authController "fita/project/coach-appointment/delivery/controller/auth"
	userController "fita/project/coach-appointment/delivery/controller/user"
	router "fita/project/coach-appointment/delivery/routers"
	"fita/project/coach-appointment/repository"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// log
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// load config
	cfg, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	dbConn := database.NewPostgresConn(cfg.DatabaseURL, cfg.DatabaseSchema)
	defer dbConn.Close()

	// initialize model
	userRepo := repository.NewUserRepo(dbConn)
	authRepo := repository.NewAuthRepo(dbConn)
	appointmentRepo := repository.NewAppointmentRepo(dbConn)
	coachAvailabilityRepo := repository.NewCoachAvailabilitymentRepo(dbConn)

	// initiate controller
	userController := userController.NewUserController(userRepo)
	authController := authController.NewAuthController(authRepo)
	appointmentController := appointmentController.NewAppointmentController(appointmentRepo, coachAvailabilityRepo)

	// create new echo
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	router.RegisterPath(e, authController, userController, appointmentController)

	// start the server, and log if it fails
	e.Logger.Fatal(e.Start(":8000"))
}
