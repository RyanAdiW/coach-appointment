package appointment

import (
	"fita/project/coach-appointment/delivery/controller"
	"fita/project/coach-appointment/delivery/middleware"
	"fita/project/coach-appointment/models"
	"fita/project/coach-appointment/repository"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	STATUS_CREATED = "CREATED"
)

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
