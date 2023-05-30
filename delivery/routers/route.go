package route

import (
	"fita/project/coach-appointment/delivery/controller/appointment"
	"fita/project/coach-appointment/delivery/controller/auth"
	"fita/project/coach-appointment/delivery/controller/user"
	"fita/project/coach-appointment/delivery/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterPath(
	e *echo.Echo,
	loginController *auth.AuthController,
	userController *user.UserController,
	appointmentController *appointment.AppointmentController) {

	// login
	e.POST("/login", loginController.AuthController())

	// user
	e.POST("/user/register", userController.CreateUserController())

	// appointment
	e.POST("/appointment/create", appointmentController.CreateAppointmentController(), middleware.JWTMiddleware())
	e.PUT("/appointment/update", appointmentController.UpdateStatusAppointment(), middleware.JWTMiddleware())
}
