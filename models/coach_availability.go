package models

type CoachAvailability struct {
	Id             int    `json:"id,omitempty"`
	UserId         string `json:"user_id,omitempty"`
	CoachName      string `json:"coach_name,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	DayOfWeek      string `json:"day_of_week,omitempty"`
	AvailableAt    string `json:"available_at,omitempty"`
	AvailableUntil string `json:"available_until,omitempty"`
}
