package controller

import "time"

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

type CreateAppointmentReq struct {
	CoachName        string    `json:"coach_name" validate:"required"`
	AppointmentStart time.Time `json:"appointment_start" validate:"required"`
	AppointmentEnd   time.Time `json:"appointment_end" validate:"required"`
}

type UpdateStatusAppointment struct {
	Id        string `json:"id" validate:"required"`
	NewStatus string `json:"new_status" validate:"required"`
}

type ReschedullingAppointment struct {
	Id               string    `json:"id" validate:"required"`
	AppointmentStart time.Time `json:"appointment_start" validate:"required"`
	AppointmentEnd   time.Time `json:"appointment_end" validate:"required"`
}

type QueryParam struct {
	Page  int `query:"page" validate:"required"`
	Limit int `query:"limit" validate:"required"`
}

type CoachAvailabilityInfo struct {
	Timezone            string `json:"timezone"`
	Day                 string `json:"day"`
	CoachAvailableFrom  string `json:"coach_available_from"`
	CoachAvailableUntil string `json:"coach_available_Until"`
	EnteredTimeFrom     string `json:"entered_time_from"`
	EnteredTimeUntil    string `json:"entered_time_until"`
}
