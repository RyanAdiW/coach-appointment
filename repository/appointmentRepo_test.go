package repository

import (
	"fita/project/coach-appointment/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateAppointment(t *testing.T) {
	Convey("Given Instance appointmentRepo", t, func() {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer db.Close()

		appointmentRepo := NewAppointmentRepo(db)

		Convey("and when CreateAppointment is called", func() {
			appointment := models.Appointment{
				UserId:           "user-1",
				Status:           "CREATED",
				CoachName:        "dipssy",
				AppointmentStart: time.Now(),
				AppointmentEnd:   time.Now(),
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			Convey("and CreateAppointment error prepare query", func() {
				query := "INSERT INTO appointment \\()"
				mock.ExpectPrepare(query)

				err := appointmentRepo.CreateAppointment(appointment)

				So(err, ShouldNotBeNil)
			})

			Convey("and CreateAppointment error exec query", func() {
				query := `INSERT INTO appointments \(user_id, status, coach_name, appointment_start, appointment_end, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(appointment.Status).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := appointmentRepo.CreateAppointment(appointment)

				So(err, ShouldNotBeNil)
			})

			Convey("and CreateAppointment success", func() {
				query := (`INSERT INTO appointments \(user_id, status, coach_name, appointment_start, appointment_end, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)`)
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(appointment.UserId, appointment.Status, appointment.CoachName, appointment.AppointmentStart, appointment.AppointmentEnd, appointment.CreatedAt, appointment.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := appointmentRepo.CreateAppointment(appointment)
				So(err, ShouldBeNil)
			})
		})
	})
}
