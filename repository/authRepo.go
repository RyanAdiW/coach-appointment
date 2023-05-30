package repository

import (
	"database/sql"
	"fmt"

	"fita/project/coach-appointment/models"
)

type AuthRepo interface {
	GetUserByEmailPass(email, password string) (*models.User, error)
	GetPasswordByEmail(email string) (string, error)
	GetIdByEmail(email string) (string, error)
	GetRole(email string) (string, error)
	GetNameByEmail(email string) (string, error)
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *authRepo {
	return &authRepo{db: db}
}

func (ar *authRepo) GetUserByEmailPass(email, password string) (*models.User, error) {
	result, err := ar.db.Query("select id, name, email, role FROM users WHERE email = $1 AND password = $2", email, password)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if isExist := result.Next(); !isExist {
		return nil, fmt.Errorf("id not found")
	}

	var user models.User
	errScan := result.Scan(&user.Id, &user.Name, &user.Email, &user.Role)
	if errScan != nil {
		return nil, errScan
	}

	return &user, nil
}

func (ar *authRepo) GetPasswordByEmail(email string) (string, error) {
	result, err := ar.db.Query("select password FROM users WHERE email = $1", email)
	if err != nil {
		return "", err
	}
	defer result.Close()
	if isExist := result.Next(); !isExist {
		return "", fmt.Errorf("id not found")
	}
	var user models.User
	errScan := result.Scan(&user.Password)
	if errScan != nil {
		return "", errScan
	}
	password := user.Password
	return password, nil
}

func (ar *authRepo) GetIdByEmail(email string) (string, error) {
	result, err := ar.db.Query("SELECT id FROM users WHERE email = $1", email)
	if err != nil {
		return "", err
	}
	defer result.Close()
	if isExist := result.Next(); !isExist {
		return "", fmt.Errorf("id not found")
	}
	var user models.User
	errScan := result.Scan(&user.Id)
	if errScan != nil {
		return "", errScan
	}
	userId := user.Id
	return userId, nil
}

func (ar *authRepo) GetRole(email string) (string, error) {
	result, err := ar.db.Query("SELECT role FROM users WHERE email = $1", email)
	if err != nil {
		return "", err
	}
	defer result.Close()
	if isExist := result.Next(); !isExist {
		return "", fmt.Errorf("id not found")
	}
	var user models.User
	errScan := result.Scan(&user.Role)
	if errScan != nil {
		return "", errScan
	}
	role := user.Role
	return role, nil
}

func (ar *authRepo) GetNameByEmail(email string) (string, error) {
	result, err := ar.db.Query("SELECT name FROM users WHERE email = $1", email)
	if err != nil {
		return "", err
	}
	defer result.Close()
	if isExist := result.Next(); !isExist {
		return "", fmt.Errorf("name not found")
	}
	var user models.User
	errScan := result.Scan(&user.Name)
	if errScan != nil {
		return "", errScan
	}
	name := user.Name
	return name, nil
}
