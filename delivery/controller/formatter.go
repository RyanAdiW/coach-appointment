package controller

type CreateUserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
}
