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

func CreateToken(userid, email, role, name string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = userid
	claims["email"] = email
	claims["role"] = role
	claims["name"] = name
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

func GetName(e echo.Context) (string, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		if name == "" {
			return name, fmt.Errorf("empty name")
		}
		return name, nil
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

func GetRole(e echo.Context) (string, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		if role == "" {
			return role, fmt.Errorf("invalid role")
		}
		return role, nil
	}
	return "", fmt.Errorf("invalid user")
}
