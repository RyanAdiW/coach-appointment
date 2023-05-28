package repository

import (
	"database/sql"
	"fita/project/coach-appointment/models"
	"log"
)

type UserRepo interface {
	CreateUser(user models.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

// create user
func (ur *userRepo) CreateUser(user models.User) error {
	query := (`INSERT INTO users (name, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`)

	statement, err := ur.db.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
