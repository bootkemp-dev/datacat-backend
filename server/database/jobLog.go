package database

import (
	"time"

	"github.com/bootkemp-dev/datacat-backend/models"
)

func (db *Database) GetJobLogsByID(jobID, limit, offset int) ([]*models.JobLog, error) {
	var logs []*models.JobLog
	rows, err := db.Query(`select * from jobLog where jobID =$1 order by id desc  limit $2 offset $3`, jobID, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var log models.JobLog

		err := rows.Scan(&log.ID, &log.UserID, &log.JobID, &log.Status, &log.LogMessage, &log.TimeChecked)
		if err != nil {
			continue
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

func (db *Database) GetAllJobLogsByID(jobID int) ([]*models.JobLog, error) {
	var logs []*models.JobLog
	rows, err := db.Query(`select * from jobLog where jobID =$1 order by id desc`, jobID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var log models.JobLog

		err := rows.Scan(&log.ID, &log.UserID, &log.JobID, &log.Status, &log.LogMessage, &log.TimeChecked)
		if err != nil {
			continue
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

func (db *Database) InsertJobLog(userID, jobID int, status, message string) error {
	stmt, err := db.Prepare(`insert into jobLog values(default, $1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID, jobID, status, message, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteJobLogs(jobID int, userID int) error {
	stmt, err := db.Prepare(`delete from jobLog where jobID = $1 and userID = $2`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(jobID, userID)
	if err != nil {
		return err
	}

	return nil
}
