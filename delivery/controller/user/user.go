package user

import (
	"fita/project/coach-appointment/delivery/controller"
	"fita/project/coach-appointment/models"
	"fita/project/coach-appointment/repository"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userRepo repository.UserRepo
}

func NewUserController(userRepo repository.UserRepo) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

func (uc *UserController) CreateUserController() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bind data
		var payload controller.CreateUserReq
		if err := c.Bind(&payload); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed binding data", err.Error()))
		}

		user := models.User{
			Name:      payload.Name,
			Email:     payload.Email,
			Password:  payload.Password,
			Role:      payload.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// create user to database
		err := uc.userRepo.CreateUser(user)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, models.InternalServerError("failed create user", err.Error()))
		}

		return c.JSON(http.StatusOK, models.SuccessOperationDefault("success", "success create user"))
	}
}
