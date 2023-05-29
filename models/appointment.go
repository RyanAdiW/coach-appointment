package models

import (
	"time"
)

type Appointment struct {
	Id               string    `json:"id"`
	UserId           string    `json:"user_id"`
	Status           string    `json:"status"`
	CoachName        string    `json:"coach_name"`
	AppointmentStart time.Time `json:"appointment_start"`
	AppointmentEnd   time.Time `json:"appointment_end"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
