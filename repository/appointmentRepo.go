package repository

import (
	"database/sql"
	"fita/project/coach-appointment/models"
	"log"
)

type AppointmentRepo interface {
	CreateAppointment(appointment models.Appointment) error
}

type appointmentRepo struct {
	db *sql.DB
}

func NewAppointmentRepo(db *sql.DB) *appointmentRepo {
	return &appointmentRepo{
		db: db,
	}
}

func (ar *appointmentRepo) CreateAppointment(appointment models.Appointment) error {
	query := (`INSERT INTO appointments (user_id, status, coach_name, appointment_start, appointment_end, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	statement, err := ar.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(appointment.UserId, appointment.Status, appointment.CoachName, appointment.AppointmentStart, appointment.AppointmentEnd, appointment.CreatedAt, appointment.UpdatedAt)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
