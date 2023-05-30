package repository

import (
	"fita/project/coach-appointment/models"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAppointmentRepo(t *testing.T) {
	Convey("Given Instance appointmentRepo", t, func() {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer db.Close()

		appointmentRepo := NewAppointmentRepo(db)

		appointment := models.Appointment{
			UserId:           "user-1",
			Status:           "CREATED",
			CoachName:        "dipssy",
			AppointmentStart: time.Now(),
			AppointmentEnd:   time.Now(),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		Convey("and when CreateAppointment is called", func() {
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

		Convey("and when UpdateStatusById is called", func() {
			Convey("and UpdateStatusById error prepare query", func() {
				query := "UPDATE appointments \\"
				mock.ExpectPrepare(query)

				err := appointmentRepo.UpdateStatusById(appointment)

				So(err, ShouldNotBeNil)
			})

			Convey("and UpdateStatusById error exec query", func() {
				query := "UPDATE appointments"
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(appointment.CoachName).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := appointmentRepo.UpdateStatusById(appointment)

				So(err, ShouldNotBeNil)
			})

			Convey("and UpdateStatusById success", func() {
				query := `UPDATE appointments SET status = \$1, updated_at = \$2 WHERE id = \$3`
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(appointment.Status, appointment.UpdatedAt, appointment.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := appointmentRepo.UpdateStatusById(appointment)

				So(err, ShouldBeNil)
			})
		})

		Convey("When GetAppointmentById is called", func() {
			appointmentID := "123"
			expectedAppointment := &models.Appointment{
				Id:               appointmentID,
				UserId:           "user-1",
				CoachName:        "dipssy",
				AppointmentStart: time.Now(),
				AppointmentEnd:   time.Now(),
			}

			Convey("and the query returns a single row", func() {
				mock.ExpectQuery("SELECT id, user_id, status, coach_name, appointment_start, appointment_end FROM appointments WHERE id =").
					WithArgs(appointmentID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "coach_name", "appointment_start", "appointment_end"}).
						AddRow(expectedAppointment.Id, expectedAppointment.UserId, expectedAppointment.Status, expectedAppointment.CoachName, expectedAppointment.AppointmentStart, expectedAppointment.AppointmentEnd))

				appointment, err := appointmentRepo.GetAppointmentById(appointmentID)

				So(err, ShouldBeNil)
				So(appointment, ShouldResemble, expectedAppointment)
			})

			Convey("and the query returns an error", func() {
				mock.ExpectQuery("SELECT id, user_id, status, coach_name, appointment_start, appointment_end FROM appointments WHERE id =").
					WithArgs(appointmentID).
					WillReturnError(fmt.Errorf("query error"))

				appointment, err := appointmentRepo.GetAppointmentById(appointmentID)

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "query error")
				So(appointment, ShouldBeNil)
			})

			Convey("and the query returns no rows", func() {
				mock.ExpectQuery("SELECT id, user_id, status, coach_name, appointment_start, appointment_end FROM appointments WHERE id =").
					WithArgs(appointmentID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "coach_name", "appointment_start", "appointment_end"}))

				appointment, err := appointmentRepo.GetAppointmentById(appointmentID)

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "sql: no rows in result set")
				So(appointment, ShouldBeNil)
			})
		})

		Convey("and when UpdateScheduleById is called", func() {
			Convey("and UpdateStatusById error prepare query", func() {
				query := "UPDATE appointments \\"
				mock.ExpectPrepare(query)

				err := appointmentRepo.UpdateScheduleById(appointment)

				So(err, ShouldNotBeNil)
			})

			Convey("and UpdateScheduleById error exec query", func() {
				query := "UPDATE appointments"
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(appointment.CoachName).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := appointmentRepo.UpdateScheduleById(appointment)

				So(err, ShouldNotBeNil)
			})

			Convey("and UpdateScheduleById success", func() {
				query := `UPDATE appointments SET appointment_start = \$1, appointment_end = \$2, updated_at = \$3, status = \$4 WHERE id = \$5`
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(appointment.AppointmentStart, appointment.AppointmentEnd, appointment.UpdatedAt, appointment.Status, appointment.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := appointmentRepo.UpdateScheduleById(appointment)

				So(err, ShouldBeNil)
			})
		})
	})
}
