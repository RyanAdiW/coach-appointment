package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte("rahasia"),
	})
}

func CreateToken(userid string, email string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = userid
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() //Token expires after 24 hours
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("rahasia"))
}

func GetEmail(e echo.Context) (string, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		if email == "" {
			return email, fmt.Errorf("empty email")
		}
		return email, nil
	}
	return "", fmt.Errorf("invalid user")
}

func GetId(e echo.Context) (string, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userid := claims["id"].(string)
		if userid == "" {
			return userid, fmt.Errorf("invalid id")
		}
		return userid, nil
	}
	return "0", fmt.Errorf("invalid user")
}

func GetIdRole(e echo.Context) (string, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		id_role := claims["id_role"].(string)
		if id_role == "" {
			return id_role, fmt.Errorf("invalid id role")
		}
		return id_role, nil
	}
	return "", fmt.Errorf("invalid user")
}
