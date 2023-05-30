package repository

import (
	"database/sql"
	"fita/project/coach-appointment/models"
	"fmt"
	"log"
)

type AppointmentRepo interface {
	CreateAppointment(appointment models.Appointment) error
	UpdateStatusById(appointment models.Appointment) error
	GetAppointmentById(id string) (*models.Appointment, error)
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
	query := (`INSERT INTO appointments (user_id, status, coach_name, appointment_start, appointment_end, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`)

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

func (ar *appointmentRepo) UpdateStatusById(appointment models.Appointment) error {
	query := (`UPDATE appointments SET status = $1 WHERE id = $2`)

	statement, err := ar.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}

	defer statement.Close()

	res, err := statement.Exec(appointment.Status, appointment.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	row, _ := res.RowsAffected()
	if row == 0 {
		return fmt.Errorf("id not found")
	}

	return nil
}

func (ar *appointmentRepo) GetAppointmentById(id string) (*models.Appointment, error) {
	var appointment models.Appointment
	row := ar.db.QueryRow(`SELECT id, user_id, status, coach_name, appointment_start, appointment_end FROM appointments WHERE id = $1`, id)

	err := row.Scan(&appointment.Id, &appointment.UserId, &appointment.Status, &appointment.CoachName, &appointment.AppointmentStart, &appointment.AppointmentEnd)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &appointment, nil
}
