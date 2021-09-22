package models

import "time"

type JobLog struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userID"`
	JobID       int       `json:"jobID"`
	Status      string    `json:"status"`
	LogMessage  string    `json:"logMessage"`
	TimeChecked time.Time `json:"timeChecked"`
}
