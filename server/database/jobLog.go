package database

import "github.com/bootkemp-dev/datacat-backend/models"

func (db *Database) GetJobLogsByID(jobID, limit, offset int) ([]*models.JobLog, error) {
	var logs []*models.JobLog
	//rows, err := db.Query(`select * from job`)

	return logs, nil
}

func (db *Database) InsertJobLog(userID, jobID int, status, message string) {
	//stmt, err := db.Prepare(`insert into jobLog values`)
}
