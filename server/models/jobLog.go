package models

import "time"

type JobLog struct {
	ID          int       `json:"id"`
	JobID       int       `json:"jobID"`
	Down        bool      `json:"down"`
	LogMessage  string    `json:"logMessage"`
	TimeChecked time.Time `json:"timeChecked"`
}
