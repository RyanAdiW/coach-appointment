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
	UpdateScheduleById(appointment models.Appointment) error
	GetAppointmentByUserId(userId string, page, limit int) ([]models.Appointment, error)
}

type appointmentRepo struct {
	db *sql.DB
}

func NewAppointmentRepo(db *sql.DB) *appointmentRepo {
	return &appointmentRepo{
		db: db,
	}
}

func (ar *appointmentRepo) CreateAppointment(req models.Appointment) error {
	query := (`INSERT INTO appointments (user_id, status, coach_name, appointment_start, appointment_end, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	statement, err := ar.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(req.UserId, req.Status, req.CoachName, req.AppointmentStart, req.AppointmentEnd, req.CreatedAt, req.UpdatedAt)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (ar *appointmentRepo) UpdateStatusById(appointment models.Appointment) error {
	query := (`UPDATE appointments SET status = $1, updated_at = $2 WHERE id = $3`)

	statement, err := ar.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}

	defer statement.Close()

	res, err := statement.Exec(appointment.Status, appointment.UpdatedAt, appointment.Id)
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

func (ar *appointmentRepo) UpdateScheduleById(appointment models.Appointment) error {
	query := (`UPDATE appointments SET appointment_start = $1, appointment_end = $2, updated_at = $3, status = $4 WHERE id = $5`)

	statement, err := ar.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}

	defer statement.Close()

	res, err := statement.Exec(appointment.AppointmentStart, appointment.AppointmentEnd, appointment.UpdatedAt, appointment.Status, appointment.Id)
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

func (ar *appointmentRepo) GetAppointmentByUserId(userId string, page, limit int) ([]models.Appointment, error) {
	var appointments []models.Appointment

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	rows, err := ar.db.Query(`SELECT id, user_id, status, coach_name, appointment_start, appointment_end, created_at, updated_at FROM appointments WHERE user_id = $1 ORDER BY updated_at DESC OFFSET $2 LIMIT $3`, userId, offset, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var appointment models.Appointment
		err = rows.Scan(&appointment.Id, &appointment.UserId, &appointment.Status, &appointment.CoachName, &appointment.AppointmentStart, &appointment.AppointmentEnd, &appointment.CreatedAt, &appointment.UpdatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}
