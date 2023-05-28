package repository

import (
	"fita/project/coach-appointment/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateUser(t *testing.T) {
	Convey("Given Instance userRepo", t, func() {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer db.Close()

		userRepository := NewUserRepo(db)

		Convey("and when CreateUser is called", func() {
			user := models.User{
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "password",
				Role:      "user",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			Convey("and CreateUser error prepare query", func() {
				query := "INSERT INTO users \\()"
				mock.ExpectPrepare(query)

				err := userRepository.CreateUser(user)

				Convey("It should return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
			Convey("and CreateUser error exec query", func() {
				query := "INSERT INTO users \\(name, email, password, role, created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)"
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(user.Name).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := userRepository.CreateUser(user)

				Convey("It should return an error", func() {
					So(err, ShouldNotBeNil)
				})
			})
			Convey("and CreateUser success", func() {
				query := "INSERT INTO users \\(name, email, password, role, created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)"
				mock.ExpectPrepare(query)
				mock.ExpectExec(query).
					WithArgs(user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := userRepository.CreateUser(user)

				Convey("It should not return an error", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
