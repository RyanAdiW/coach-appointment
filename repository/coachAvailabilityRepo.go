package repository

import (
	"database/sql"
	"fita/project/coach-appointment/models"
	"fmt"
)

type CoachAvailabilityRepo interface {
	GetAvailability(coachName string) (*models.CoachAvailability, error)
}

type coachAvailabilityRepo struct {
	db *sql.DB
}

func NewCoachAvailabilitymentRepo(db *sql.DB) *coachAvailabilityRepo {
	return &coachAvailabilityRepo{
		db: db,
	}
}

func (cr *coachAvailabilityRepo) GetAvailability(coachName string) (*models.CoachAvailability, error) {
	result, err := cr.db.Query("SELECT user_id, coach_name, timezone, day_of_week, available_at, available_until FROM coach_availability WHERE coach_name = $1", coachName)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	if isExist := result.Next(); !isExist {
		return nil, fmt.Errorf("coach not found")
	}
	var coachAvailability models.CoachAvailability
	errScan := result.Scan(&coachAvailability.UserId, &coachAvailability.CoachName, &coachAvailability.Timezone, &coachAvailability.DayOfWeek, &coachAvailability.AvailableAt, &coachAvailability.AvailableUntil)
	if errScan != nil {
		return nil, errScan
	}

	return &coachAvailability, nil
}
