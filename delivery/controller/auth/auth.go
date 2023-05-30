package auth

import (
	"fita/project/coach-appointment/delivery/controller"
	"fita/project/coach-appointment/delivery/middleware"
	"fita/project/coach-appointment/models"
	"fita/project/coach-appointment/repository"

	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authRepo repository.AuthRepo
}

func NewAuthController(authRepo repository.AuthRepo) *AuthController {
	return &AuthController{
		authRepo: authRepo,
	}
}

func (ac *AuthController) AuthController() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload controller.LoginReq

		//bind request data
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.BadRequest("unauthorized", "failed to bind"))
		}

		passwordOnDb, err := ac.authRepo.GetPasswordByEmail(payload.Email)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadRequest("unauthorized", "email not found"))
		}

		if payload.Password != passwordOnDb {
			return c.JSON(http.StatusBadRequest, models.BadRequest("unauthorized", "invalid password"))
		}

		// get token from login credential
		user, err := ac.authRepo.GetUserByEmailPass(payload.Email, payload.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadRequest("unauthorized", "failed to create token"))
		}

		token, _ := middleware.CreateToken(user.Id, user.Email, user.Role, user.Name)

		uid, _ := ac.authRepo.GetIdByEmail(payload.Email)
		role, _ := ac.authRepo.GetRole(payload.Email)
		name, _ := ac.authRepo.GetNameByEmail(payload.Email)

		data := controller.LoginResponse{
			Token:  token,
			UserId: uid,
			Role:   role,
			Name:   name,
		}
		return c.JSON(http.StatusOK, models.SuccessOperationWithData("success", "login success", data))
	}
}
