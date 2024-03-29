package repository

import (
	"errors"
	"fita/project/coach-appointment/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoachAvailabilityRepo(t *testing.T) {
	Convey("Given an instance of coachAvailabilityRepo", t, func() {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer db.Close()

		coachAvailabilityRepo := NewCoachAvailabilitymentRepo(db)

		Convey("and when GetAvailability is called", func() {
			coachName := "dipssy"
			dayOfWeek := "Monday"
			expectedCoachAvailability := &models.CoachAvailability{
				UserId:         "user-1",
				CoachName:      "dipssy",
				Timezone:       "UTC",
				DayOfWeek:      "Monday",
				AvailableAt:    "08:00",
				AvailableUntil: "16:00",
			}

			Convey("and error coach is not found in the database", func() {
				rows := sqlmock.NewRows([]string{})

				query := "SELECT user_id, coach_name, timezone, day_of_week, available_at, available_until FROM coach_availability WHERE coach_name = $1 and day_of_week = $2"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(coachName, dayOfWeek).WillReturnRows(rows)

				coachAvailability, err := coachAvailabilityRepo.GetAvailability(coachName, dayOfWeek)

				So(err, ShouldNotBeNil)
				So(coachAvailability, ShouldBeNil)
			})

			Convey("and an error occurs during query execution", func() {
				expectedErr := errors.New("database error")

				query := "SELECT user_id, coach_name, timezone, day_of_week, available_at, available_until FROM coach_availability WHERE coach_name = $1 and day_of_week = $2"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(coachName, dayOfWeek).WillReturnError(expectedErr)

				coachAvailability, err := coachAvailabilityRepo.GetAvailability(coachName, dayOfWeek)

				So(err, ShouldEqual, expectedErr)
				So(coachAvailability, ShouldBeNil)
			})

			Convey("and GetAvailability success", func() {
				rows := sqlmock.NewRows([]string{"user_id", "coach_name", "timezone", "day_of_week", "available_at", "available_until"}).
					AddRow(expectedCoachAvailability.UserId, expectedCoachAvailability.CoachName, expectedCoachAvailability.Timezone, expectedCoachAvailability.DayOfWeek, expectedCoachAvailability.AvailableAt, expectedCoachAvailability.AvailableUntil)

				query := "SELECT user_id, coach_name, timezone, day_of_week, available_at, available_until FROM coach_availability WHERE coach_name = $1 and day_of_week = $2"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(coachName, dayOfWeek).WillReturnRows(rows)

				coachAvailability, err := coachAvailabilityRepo.GetAvailability(coachName, dayOfWeek)

				So(err, ShouldBeNil)
				So(coachAvailability, ShouldResemble, expectedCoachAvailability)
			})
		})

		Convey("and when GetCoachTimezone is called", func() {
			coachName := "dipssy"
			expectedCoachTimezone := "UTC"

			Convey("and error coachTimezone is not found in the database", func() {
				rows := sqlmock.NewRows([]string{})

				query := "SELECT timezone WHERE coach_name = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(coachName).WillReturnRows(rows)

				timezone, err := coachAvailabilityRepo.GetCoachTimezone(coachName)

				So(err, ShouldNotBeNil)
				So(timezone, ShouldEqual, "")
			})

			Convey("and an error occurs during query execution", func() {
				expectedErr := errors.New("database error")

				query := "SELECT timezone FROM coach_availability WHERE coach_name = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(coachName).WillReturnError(expectedErr)

				coachAvailability, err := coachAvailabilityRepo.GetCoachTimezone(coachName)

				So(err, ShouldEqual, expectedErr)
				So(coachAvailability, ShouldEqual, "")
			})

			Convey("and GetAvailability success", func() {
				rows := sqlmock.NewRows([]string{"timezone"}).
					AddRow(expectedCoachTimezone)

				query := "SELECT timezone FROM coach_availability WHERE coach_name = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(coachName).WillReturnRows(rows)

				coachAvailability, err := coachAvailabilityRepo.GetCoachTimezone(coachName)

				So(err, ShouldBeNil)
				So(coachAvailability, ShouldResemble, expectedCoachTimezone)
			})
		})
	})
}
