package models

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type NewJobRequest struct {
	JobName   string `json:"name"`
	JobURL    string `json:"url"`
	Frequency int64  `json:"frequency"`
}

func NewPool() Pool {
	return Pool{Jobs: []*Job{}}
}

type Pool struct {
	Jobs []*Job
}

func (p Pool) GetJob(jobID int, userID int) (*Job, error) {
	for _, v := range p.Jobs {
		if v.UserID == userID && v.ID == jobID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Job not found")
}

func (p Pool) RemoveJob(jobID int, userID int) error {
	for i := range p.Jobs {
		if p.Jobs[i].ID == jobID && p.Jobs[i].UserID == userID {
			p.Jobs[i].Stop()
			copy(p.Jobs[i:], p.Jobs[i+1:])
			p.Jobs[len(p.Jobs)-1] = nil
			p.Jobs = p.Jobs[:len(p.Jobs)-1]
			return nil
		}
	}
	return fmt.Errorf("job not found")
}

type Job struct {
	ID         int       `json:"id"`
	Name       string    `json:"job_name"`
	URL        string    `json:"job_url"`
	Frequency  int64     `json:"frequency"`
	UserID     int       `json:"user_id"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	status     string
	done       chan bool
}

func NewJob(jobId int, userID int, name, url string, freq int64) Job {
	j := Job{
		ID:         jobId,
		Name:       name,
		URL:        url,
		Frequency:  freq,
		UserID:     userID,
		Active:     false,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		status:     "NA",
		done:       make(chan bool),
	}

	return j
}

func (j Job) Run() {
	j.Active = true
	log.Printf("Starting job | ID: %d | Name: %s\n", j.ID, j.Name)
	go func() {
		for {
			select {
			case <-j.done:
				j.Active = false
				log.Printf("Ending job | ID: %d | Name: %s | URL: %s\n", j.ID, j.Name, j.URL)
				return
			default:
				log.Printf("Job with ID: %d checking status of: %s ", j.ID, j.URL)
				err := j.URLStatus()
				if err != nil {
					j.status = "DOWN"
				}

				time.Sleep(time.Duration(j.Frequency))
			}
		}
	}()
}

func (j Job) URLStatus() error {
	resp, err := http.Get(j.URL)
	if err == nil && resp.StatusCode == 200 {
		return nil
	} else {
		return err
	}
}

func (j Job) Stop() error {
	if j.Active == false {
		return fmt.Errorf("Job is not running")
	} else {
		j.done <- true
	}

	return nil
}

func (j Job) GetStatus() string {
	return j.status
}
