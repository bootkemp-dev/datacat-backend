package models

type Job struct {
	ID        int    `json:"id"`
	JobName   string `json:"job_name"`
	JobURL    string `json:"job_url"`
	Frequency int    `json:"frequency"`
	UserID    int    `json:"user_id"`
}

type NewJobRequest struct {
	JobName   string `json:"name"`
	JobURL    string `json:"url"`
	Frequency int    `json:"frequency"`
}
