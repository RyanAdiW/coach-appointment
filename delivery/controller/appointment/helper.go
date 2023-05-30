package appointment

import (
	"fita/project/coach-appointment/delivery/controller"
	"fita/project/coach-appointment/delivery/middleware"
	"fita/project/coach-appointment/models"
	"fita/project/coach-appointment/repository"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func UpdateStatusByCoach(c echo.Context, appointmentRepo repository.AppointmentRepo, payload controller.UpdateStatusAppointment) error {
	role, err := middleware.GetRole(c)
	if err != nil {
		return fmt.Errorf("GetRole error")
	}

	if strings.ToUpper(role) != ROLE_COACH {
		return fmt.Errorf("role unauthorized")
	}

	appointment, _ := appointmentRepo.GetAppointmentById(payload.Id)
	if appointment == nil {
		return fmt.Errorf("appointment not found")
	}

	name, err := middleware.GetName(c)
	if err != nil {
		return fmt.Errorf("GetName error")
	}

	if appointment.CoachName != name {
		return fmt.Errorf("name unauthorized")
	}

	switch payload.NewStatus {
	case STATUS_ACCEPTED:
		if strings.ToUpper(appointment.Status) != STATUS_CREATED && strings.ToUpper(appointment.Status) != STATUS_RESCHEDULING {
			return fmt.Errorf("current status must be CREATED or RESCHEDULING")
		}

		updateStatReq := models.Appointment{
			Id:        payload.Id,
			Status:    STATUS_ACCEPTED,
			UpdatedAt: time.Now(),
		}

		err = appointmentRepo.UpdateStatusById(updateStatReq)
		if err != nil {
			return fmt.Errorf("failed update appointment status to COACH_ACCEPTED, error from db")
		}
	case STATUS_REJECTED:
		if strings.ToUpper(appointment.Status) != STATUS_CREATED && strings.ToUpper(appointment.Status) != STATUS_RESCHEDULING {
			return fmt.Errorf("current status must be CREATED or RESCHEDULING")
		}

		updateStatReq := models.Appointment{
			Id:        payload.Id,
			Status:    STATUS_REJECTED,
			UpdatedAt: time.Now(),
		}

		err = appointmentRepo.UpdateStatusById(updateStatReq)
		if err != nil {
			return fmt.Errorf("failed update appointment status to COACH_REJECTED, error from db")
		}
	case STATUS_RESCHEDULE_REQUESTED:
		if strings.ToUpper(appointment.Status) != STATUS_CREATED {
			return fmt.Errorf("current status must be CREATED")
		}

		updateStatReq := models.Appointment{
			Id:        payload.Id,
			Status:    STATUS_RESCHEDULE_REQUESTED,
			UpdatedAt: time.Now(),
		}

		err = appointmentRepo.UpdateStatusById(updateStatReq)
		if err != nil {
			return fmt.Errorf("failed update appointment status to RESCHEDULE_REQUESTED, error from db")
		}
	}
	return nil
}

func UpdateStatusByUser(c echo.Context, appointmentRepo repository.AppointmentRepo, payload controller.UpdateStatusAppointment) error {
	role, err := middleware.GetRole(c)
	if err != nil {
		return fmt.Errorf("GetRole error")
	}

	if strings.ToUpper(role) != ROLE_USER {
		return fmt.Errorf("role unauthorized")
	}

	appointment, _ := appointmentRepo.GetAppointmentById(payload.Id)
	if appointment == nil {
		return fmt.Errorf("appointment not found")
	}

	userId, err := middleware.GetId(c)
	if err != nil {
		return fmt.Errorf("GetId error")
	}

	if appointment.UserId != userId {
		return fmt.Errorf("userID unauthorized")
	}

	switch payload.NewStatus {
	case STATUS_RESCHEDULE_REJECTED:
		if strings.ToUpper(appointment.Status) != STATUS_RESCHEDULE_REQUESTED {
			return fmt.Errorf("current status must be RESCHEDULE_REQUESTED")
		}

		updateStatReq := models.Appointment{
			Id:        payload.Id,
			Status:    STATUS_RESCHEDULE_REJECTED,
			UpdatedAt: time.Now(),
		}

		err = appointmentRepo.UpdateStatusById(updateStatReq)
		if err != nil {
			return fmt.Errorf("failed update appointment status to RESCHEDULE_REJECTED, error from db")
		}
	}
	return nil
}

func ReschedullingByUser(c echo.Context,
	appointment *models.Appointment,
	appointmentRepo repository.AppointmentRepo,
	payload controller.ReschedullingAppointment) error {
	role, err := middleware.GetRole(c)
	if err != nil {
		return fmt.Errorf("GetRole error")
	}

	if strings.ToUpper(role) != ROLE_USER {
		return fmt.Errorf("role unauthorized")
	}

	userId, err := middleware.GetId(c)
	if err != nil {
		return fmt.Errorf("GetId error")
	}

	if appointment.UserId != userId {
		return fmt.Errorf("userID unauthorized")
	}

	if strings.ToUpper(appointment.Status) != STATUS_RESCHEDULE_REQUESTED {
		return fmt.Errorf("current status must be RESCHEDULE_REQUESTED")
	}

	req := models.Appointment{
		Id:               payload.Id,
		Status:           STATUS_RESCHEDULING,
		AppointmentStart: payload.AppointmentStart,
		AppointmentEnd:   payload.AppointmentEnd,
		UpdatedAt:        time.Now(),
	}

	err = appointmentRepo.UpdateScheduleById(req)
	if err != nil {
		return fmt.Errorf("failed RESCHEDULLING, error from db")
	}

	return nil
}
