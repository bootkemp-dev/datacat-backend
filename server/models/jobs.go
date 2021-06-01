package models

import "time"

type Job struct {
	ID        int       `json:"id"`
	JobName   string    `json:"job_name"`
	JobURL    string    `json:"job_url"`
	Frequency uint      `json:"frequency"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
type NewJobRequest struct {
	JobName   string `json:"name"`
	JobURL    string `json:"url"`
	Frequency uint   `json:"frequency"`
}
