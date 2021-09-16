package database

import "github.com/bootkemp-dev/datacat-backend/models"

func (db *Database) GetJobLogsByID(jobID int) ([]*models.JobLog, error) {
	var logs []*models.JobLog
	//rows, err := db.Query(`select * from job`)

	return logs, nil
}

func (db *Database) InsertJobLog() {}
