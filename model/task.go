package model

import "time"

// Task if the value of any entity is default value (0,"",false) than omitempty will not
// include the field in json response.
type Task struct {
	Id            int64     `json:"id"`
	UserId        int64     `json:"user_id,omitempty"`
	Description   string    `json:"description"`
	StartedAt     time.Time `json:"started_at"`
	EndedAt       time.Time `json:"ended_at"`
	TotalDuration int       `json:"total_duration"`
}
