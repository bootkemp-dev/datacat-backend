package models

import "time"

type Job struct {
	ID         int       `json:"id"`
	Name       string    `json:"job_name"`
	URL        string    `json:"job_url"`
	Frequency  uint      `json:"frequency"`
	UserID     int       `json:"user_id"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Status     string    `json:"status"`
}

type NewJobRequest struct {
	JobName   string `json:"name"`
	JobURL    string `json:"url"`
	Frequency int64  `json:"frequency"`
}
