package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
)

func createTestContextWithToken(token *jwt.Token) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", token)
	return c
}

func TestGetEmail(t *testing.T) {
	Convey("Given TestGetEmail Instance", t, func() {
		Convey("GetEmail success", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "123",
				"email":      "user@example.com",
				"id_role":    "admin",
				"exp":        time.Now().Add(time.Hour * 24).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token.Valid = true
			c := createTestContextWithToken(token)

			email, err := GetEmail(c)

			So(err, ShouldBeNil)
			So(email, ShouldEqual, "user@example.com")
		})
		Convey("and GetEmail error empty email", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "123",
				"email":      "",
				"id_role":    "admin",
				"exp":        time.Now().Add(time.Hour * 24).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token.Valid = true
			c := createTestContextWithToken(token)

			_, err := GetEmail(c)

			So(err, ShouldNotBeNil)
		})
		Convey("and GetEmail error user invalid", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "123",
				"email":      "user@example.com",
				"id_role":    "admin",
				"exp":        time.Now().Add(time.Hour * 24).Unix(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			c := createTestContextWithToken(token)

			_, err := GetEmail(c)

			So(err, ShouldNotBeNil)
		})
	})

}

func TestGetId(t *testing.T) {
	Convey("Given GetId Instance", t, func() {
		Convey("and GetId success", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "123",
				"email":      "user@example.com",
				"id_role":    "admin",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token.Valid = true

			c := createTestContextWithToken(token)

			id, err := GetId(c)

			So(err, ShouldBeNil)
			So(id, ShouldEqual, "123")
		})
		Convey("and GetId error empty user id", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "",
				"email":      "user@example.com",
				"id_role":    "admin",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token.Valid = true

			c := createTestContextWithToken(token)

			_, err := GetId(c)

			So(err, ShouldNotBeNil)
		})
		Convey("and GetId error user invalid", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "123",
				"email":      "user@example.com",
				"id_role":    "admin",
				"exp":        time.Now().Add(time.Hour * 24).Unix(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			c := createTestContextWithToken(token)

			_, err := GetId(c)

			So(err, ShouldNotBeNil)
		})
	})

}

func TestGetIdRole(t *testing.T) {
	Convey("Given Instance GetRole", t, func() {
		Convey("and GetRole success", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         123,
				"email":      "user@example.com",
				"role":       "coach",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token.Valid = true

			c := createTestContextWithToken(token)

			idRole, err := GetRole(c)

			So(err, ShouldBeNil)
			So(idRole, ShouldEqual, "coach")
		})
		Convey("and GetRole error empty idRole", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "",
				"email":      "user@example.com",
				"role":       "",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token.Valid = true

			c := createTestContextWithToken(token)

			_, err := GetRole(c)

			So(err, ShouldNotBeNil)
		})
		Convey("and GetRole error user invalid", func() {
			claims := jwt.MapClaims{
				"authorized": true,
				"id":         "123",
				"email":      "user@example.com",
				"role":       "admin",
				"exp":        time.Now().Add(time.Hour * 24).Unix(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			c := createTestContextWithToken(token)

			_, err := GetRole(c)

			So(err, ShouldNotBeNil)
		})
	})

}

func TestCreateToken(t *testing.T) {
	Convey("Given CrateToken Instances", t, func() {
		Convey("and createToken success", func() {
			userID := "123"
			email := "user@example.com"
			role := "coach"

			tokenString, err := CreateToken(userID, email, role)
			So(err, ShouldBeNil)
			So(tokenString, ShouldNotEqual, "")

			// Parse the token to validate the claims
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte("rahasia"), nil
			})
			So(err, ShouldBeNil)
			So(token.Valid, ShouldEqual, true)

			// Validate the claims
			claims := token.Claims.(jwt.MapClaims)
			So(claims["authorized"].(bool), ShouldEqual, true)
			So(claims["id"].(string), ShouldEqual, userID)
			So(claims["email"].(string), ShouldEqual, email)
			So(claims["role"].(string), ShouldEqual, role)
		})
	})
}
