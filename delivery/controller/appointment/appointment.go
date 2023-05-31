package appointment

import (
	"fita/project/coach-appointment/delivery/controller"
	"fita/project/coach-appointment/delivery/middleware"
	"fita/project/coach-appointment/models"
	"fita/project/coach-appointment/repository"
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

		// validate coach available time
		// ===========================
		timezone, err := ac.coachAvailabilityRepo.GetCoachTimezone(payload.CoachName)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("get coach timezone failed", err.Error()))
		}

		// convert payload local time to coach local time
		location, err := time.LoadLocation(timezone)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("failed LoadLocation time", err.Error()))
		}

		convertedAppointmentStart := payload.AppointmentStart.In(location)
		convertedAppointmentEnd := payload.AppointmentEnd.In(location)

		//get coach avaiability
		coachAvlb, err := ac.coachAvailabilityRepo.GetAvailability(payload.CoachName, convertedAppointmentStart.Weekday().String())
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("get coach avlb failed", err.Error()))
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
		// ===========================

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
			log.Println(err)
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

		payload.NewStatus = strings.ToUpper(payload.NewStatus)
		if payload.NewStatus == STATUS_ACCEPTED || payload.NewStatus == STATUS_REJECTED || payload.NewStatus == STATUS_RESCHEDULE_REQUESTED {
			err := UpdateStatusByCoach(c, ac.appointmentRepo, payload)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusBadRequest, models.BadRequest("failed", err.Error()))
			}
		} else if payload.NewStatus == STATUS_RESCHEDULE_REJECTED {
			err := UpdateStatusByUser(c, ac.appointmentRepo, payload)
			if err != nil {
				log.Println(err)
				return c.JSON(http.StatusBadRequest, models.BadRequest("failed", err.Error()))
			}
		}

		return c.JSON(http.StatusOK, models.SuccessOperationDefault("success", "success update status appointment"))
	}
}

func (ac *AppointmentController) ReschedullingByUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bind data
		var payload controller.ReschedullingAppointment
		if err := c.Bind(&payload); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed binding data", err.Error()))
		}

		appointment, _ := ac.appointmentRepo.GetAppointmentById(payload.Id)
		if appointment == nil {
			return c.JSON(http.StatusNotFound, models.NotFound("failed", "appointment not found"))
		}

		// validate coach available time
		// ===========================
		timezone, err := ac.coachAvailabilityRepo.GetCoachTimezone(appointment.CoachName)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("get coach timezone failed", err.Error()))
		}

		// convert payload local time to coach local time
		location, err := time.LoadLocation(timezone)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("failed LoadLocation time", err.Error()))
		}

		convertedAppointmentStart := payload.AppointmentStart.In(location)
		convertedAppointmentEnd := payload.AppointmentEnd.In(location)

		//get coach avaiability
		coachAvlb, err := ac.coachAvailabilityRepo.GetAvailability(appointment.CoachName, convertedAppointmentStart.Weekday().String())
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, models.InternalServerError("get coach avlb failed", err.Error()))
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
		// ===========================

		err = ReschedullingByUser(c, appointment, ac.appointmentRepo, payload)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed", err.Error()))
		}

		return c.JSON(http.StatusOK, models.SuccessOperationDefault("success", "success update status appointment"))
	}
}

func (ac *AppointmentController) GetAppointmentByUserId() echo.HandlerFunc {
	return func(c echo.Context) error {
		// bind data
		var payload controller.QueryParamGetAppointmentByUserId
		if err := c.Bind(&payload); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, models.BadRequest("failed binding data", err.Error()))
		}

		userId, err := middleware.GetId(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, models.InternalServerError("failed", "get userId error"))
		}

		appointments, err := ac.appointmentRepo.GetAppointmentByUserId(userId, payload.Page, payload.Limit)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, models.InternalServerError("failed", "get appointments error"))
		}

		return c.JSON(http.StatusOK, models.SuccessOperationWithData("success", "success get appointments", appointments))
	}
}
