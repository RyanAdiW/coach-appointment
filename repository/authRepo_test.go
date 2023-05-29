package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthRepo(t *testing.T) {
	Convey("Given an instance of authRepo", t, func() {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer db.Close()

		authRepo := NewAuthRepo(db)

		Convey("When LoginEmail is called", func() {
			email := "test@example.com"
			password := "password123"

			Convey("and the email and password are valid", func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
					AddRow("user-1", "John Doe", email, "user")

				query := "select id, name, email, role FROM users WHERE email = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email, password).WillReturnRows(rows)

				mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 0)) // Ignore middleware token creation

				_, err := authRepo.LoginEmail(email, password)

				So(err, ShouldBeNil)
			})

			Convey("and the email and password are invalid", func() {
				rows := sqlmock.NewRows([]string{})

				query := "select id, name, email, role FROM users WHERE email = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email, password).WillReturnRows(rows)

				token, err := authRepo.LoginEmail(email, password)

				So(err.Error(), ShouldEqual, "id not found")
				So(token, ShouldEqual, "")
			})

			Convey("and an error occurs during query execution", func() {
				expectedErr := errors.New("database error")

				query := "select id, name, email, role FROM users WHERE email = $1 AND password = $2"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email, password).WillReturnError(expectedErr)

				token, err := authRepo.LoginEmail(email, password)

				So(err, ShouldEqual, expectedErr)
				So(token, ShouldEqual, "")
			})
		})

		Convey("When GetPasswordByEmail is called", func() {
			email := "test@example.com"
			expectedPassword := "password123"

			Convey("and the email is valid", func() {
				rows := sqlmock.NewRows([]string{"password"}).
					AddRow(expectedPassword)

				query := "select password FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				password, err := authRepo.GetPasswordByEmail(email)

				So(err, ShouldBeNil)
				So(password, ShouldEqual, expectedPassword)
			})

			Convey("and the email is invalid", func() {
				rows := sqlmock.NewRows([]string{})

				query := "select password FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				password, err := authRepo.GetPasswordByEmail(email)

				So(err.Error(), ShouldEqual, "id not found")
				So(password, ShouldEqual, "")
			})

			Convey("and an error occurs during query execution", func() {
				expectedErr := errors.New("database error")

				query := "select password FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnError(expectedErr)

				password, err := authRepo.GetPasswordByEmail(email)

				So(err, ShouldEqual, expectedErr)
				So(password, ShouldEqual, "")
			})
		})

		Convey("When GetIdByEmail is called", func() {
			email := "test@example.com"
			expectedID := "user-1"

			Convey("and the email is valid", func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(expectedID)

				query := "SELECT id FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				id, err := authRepo.GetIdByEmail(email)

				So(err, ShouldBeNil)
				So(id, ShouldEqual, expectedID)
			})

			Convey("and the email is invalid", func() {
				rows := sqlmock.NewRows([]string{})

				query := "SELECT id FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				id, err := authRepo.GetIdByEmail(email)

				So(err.Error(), ShouldEqual, "id not found")
				So(id, ShouldEqual, "")
			})

			Convey("and an error occurs during query execution", func() {
				expectedErr := errors.New("database error")

				query := "SELECT id FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnError(expectedErr)

				id, err := authRepo.GetIdByEmail(email)

				So(err, ShouldEqual, expectedErr)
				So(id, ShouldEqual, "")
			})
		})

		Convey("When GetRole is called", func() {
			email := "test@example.com"
			expectedRole := "user"

			Convey("and the email is valid", func() {
				rows := sqlmock.NewRows([]string{"role"}).
					AddRow(expectedRole)

				query := "SELECT role FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				role, err := authRepo.GetRole(email)

				So(err, ShouldBeNil)
				So(role, ShouldEqual, expectedRole)
			})

			Convey("and the email is invalid", func() {
				rows := sqlmock.NewRows([]string{})

				query := "SELECT role FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				role, err := authRepo.GetRole(email)

				So(err.Error(), ShouldEqual, "id not found")
				So(role, ShouldEqual, "")
			})

			Convey("and an error occurs during query execution", func() {
				expectedErr := errors.New("database error")

				query := "SELECT role FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnError(expectedErr)

				role, err := authRepo.GetRole(email)

				So(err, ShouldEqual, expectedErr)
				So(role, ShouldEqual, "")
			})
		})

		Convey("When GetNameByEmail is called", func() {
			email := "test@example.com"
			expectedName := "John Doe"

			Convey("nd the email is valid", func() {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow(expectedName)

				query := "SELECT name FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				name, err := authRepo.GetNameByEmail(email)

				So(err, ShouldBeNil)
				So(name, ShouldEqual, expectedName)
			})

			Convey("and the email is invalid", func() {
				rows := sqlmock.NewRows([]string{})

				query := "SELECT name FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

				name, err := authRepo.GetNameByEmail(email)

				So(err.Error(), ShouldEqual, "name not found")
				So(name, ShouldEqual, "")
			})

			Convey("and an error occurs during query execution", func() {
				expectedErr := errors.New("database error")

				query := "SELECT name FROM users WHERE email = $1"
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnError(expectedErr)

				name, err := authRepo.GetNameByEmail(email)

				So(err, ShouldEqual, expectedErr)
				So(name, ShouldEqual, "")
			})
		})

	})
}
