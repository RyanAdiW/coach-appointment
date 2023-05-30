package appointment

import (
	"fita/project/coach-appointment/delivery/controller"
	"fita/project/coach-appointment/delivery/middleware"
	"fita/project/coach-appointment/models"
	"fita/project/coach-appointment/repository"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	STATUS_CREATED              = "CREATED"              // by user
	STATUS_ACCEPTED             = "COACH_ACCEPTED"       // by coach
	STATUS_REJECTED             = "COACH_REJECTED"       // by coach
	STATUS_RESCHEDULE_REQUESTED = "RESCHEDULE_REQUESTED" // by coach
	STATUS_RESCHEDULE_REJECTED  = "RESCHEDULE_REJECTED"  // by user
	STATUS_RESCHEDULING         = "RESCHEDULING"         // by user

	ROLE_COACH = "COACH"
	ROLE_USER  = "USER"
)

var statusMap = map[string]bool{
	STATUS_CREATED:              true,
	STATUS_ACCEPTED:             true,
	STATUS_REJECTED:             true,
	STATUS_RESCHEDULE_REQUESTED: true,
	STATUS_RESCHEDULE_REJECTED:  true,
	STATUS_RESCHEDULING:         true,
}

type AppointmentController struct {
	appointmentRepo       repository.AppointmentRepo
	coachAvailabilityRepo repository.CoachAvailabilityRepo
}

func NewAppointmentController(appointmentRepo repository.AppointmentRepo, coachAvailabilityRepo repository.CoachAvailabilityRepo) *AppointmentController {
	return &AppointmentController{
		appointmentRepo:       appointmentRepo,
		coachAvailabilityRepo: coachAvailabilityRepo,
	}
}

func (ac *AppointmentController) CreateAppointmentController() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bind data
		var payload controller.CreateAppointmentReq
		if err := c.Bind(&payload); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed binding data", err.Error()))
		}

		//get coach avaiability
		coachAvlb, err := ac.coachAvailabilityRepo.GetAvailability(payload.CoachName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("get coach avlb failed", err.Error()))
		}

		// convert payload local time to coach local time
		location, err := time.LoadLocation(coachAvlb.Timezone)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("failed LoadLocation time", err.Error()))
		}

		convertedAppointmentStart := payload.AppointmentStart.In(location)
		convertedAppointmentEnd := payload.AppointmentEnd.In(location)

		if coachAvlb.DayOfWeek != convertedAppointmentStart.Weekday().String() {
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed", "coach not available on selected day"))
		}

		avlbAtConv, _ := time.Parse("15:04:05", coachAvlb.AvailableAt)
		avlbUntilConv, _ := time.Parse("15:04:05", coachAvlb.AvailableUntil)

		payloadStart, _ := time.Parse("15:04:05", convertedAppointmentStart.Format("15:04:05"))
		payloadEnd, _ := time.Parse("15:04:05", convertedAppointmentEnd.Format("15:04:05"))

		isInvalidTime := payloadStart.Before(avlbAtConv) ||
			payloadStart.After(avlbUntilConv) ||
			payloadEnd.Before(avlbAtConv) ||
			payloadEnd.After(avlbUntilConv)

		if isInvalidTime {
			data := controller.CoachAvailabilityInfo{
				Timezone:            coachAvlb.Timezone,
				Day:                 coachAvlb.DayOfWeek,
				CoachAvailableFrom:  coachAvlb.AvailableAt,
				CoachAvailableUntil: coachAvlb.AvailableUntil,
				EnteredTimeFrom:     convertedAppointmentStart.Format("15:04:05"),
				EnteredTimeUntil:    convertedAppointmentEnd.Format("15:04:05"),
			}
			return c.JSON(http.StatusBadRequest, models.BadRequestWithData("failed", "coach not available, choose valid time:", data))
		}

		// set userId from jwt
		userId, _ := middleware.GetId(c)

		insertDb := models.Appointment{
			UserId:           userId,
			Status:           STATUS_CREATED,
			CoachName:        coachAvlb.CoachName,
			AppointmentStart: convertedAppointmentStart,
			AppointmentEnd:   convertedAppointmentEnd,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		err = ac.appointmentRepo.CreateAppointment(insertDb)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("failed", "failed insert request to db"))
		}

		return c.JSON(http.StatusOK, models.SuccessOperationDefault("success", "success create appointment"))
	}
}

func (ac *AppointmentController) UpdateStatusAppointment() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bind data
		var payload controller.UpdateStatusAppointment
		if err := c.Bind(&payload); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed binding data", err.Error()))
		}

		if !statusMap[strings.ToUpper(payload.NewStatus)] {
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed", "status not allowed"))
		}

		switch strings.ToUpper(payload.NewStatus) {
		case STATUS_ACCEPTED:
			err := AccetpStatus(c, ac.appointmentRepo, payload)
			if err != nil {
				return c.JSON(http.StatusBadRequest, models.BadRequest("failed", err.Error()))
			}
		}

		return c.JSON(http.StatusOK, models.SuccessOperationDefault("success", "success update status appointment"))
	}
}

func AccetpStatus(c echo.Context, appointmentRepo repository.AppointmentRepo, payload controller.UpdateStatusAppointment) error {
	role, err := middleware.GetRole(c)
	if err != nil {
		return fmt.Errorf("GetRole error")
	}

	if strings.ToUpper(role) != ROLE_COACH {
		return fmt.Errorf("only COACH can ACCEPT request")
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
		return fmt.Errorf("Unauthorized")
	}

	if strings.ToUpper(appointment.Status) != STATUS_CREATED && strings.ToUpper(appointment.Status) != STATUS_RESCHEDULING {
		return fmt.Errorf("current status must be CREATED or RESCHEDULING")
	}

	updateStatReq := models.Appointment{
		Id:     payload.Id,
		Status: STATUS_ACCEPTED,
	}

	err = appointmentRepo.UpdateStatusById(updateStatReq)
	if err != nil {
		return fmt.Errorf("failed update appointment status to COACH_ACCEPTED")
	}

	return nil
}
