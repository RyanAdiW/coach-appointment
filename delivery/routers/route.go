package route

import (
	"fita/project/coach-appointment/delivery/controller/auth"
	"fita/project/coach-appointment/delivery/controller/user"

	"github.com/labstack/echo/v4"
)

func RegisterPath(
	e *echo.Echo,
	loginController *auth.AuthController,
	userController *user.UserController) {

	// login
	e.POST("/login", loginController.AuthController())

	// user
	e.POST("/user/register", userController.CreateUserController())
}
