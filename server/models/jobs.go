package models

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bootkemp-dev/datacat-backend/logger"
)

type NewJobRequest struct {
	JobName   string `json:"name"`
	JobURL    string `json:"url"`
	Frequency int64  `json:"frequency"`
}

func NewPool() Pool {
	return Pool{jobs: []*Job{}}
}

type Pool struct {
	jobs []*Job
}

func (p *Pool) AddJob(job *Job) {
	p.jobs = append(p.jobs, job)
}

func (p *Pool) RemoveJob(jobID int, userID int) error {
	for i := range p.jobs {
		if p.jobs[i].ID == jobID && p.jobs[i].UserID == userID {
			if p.jobs[i].Active == true {
				go p.jobs[i].Stop()
			}
			p.jobs = append(p.jobs[:i], p.jobs[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Job not found in the pool")
}

func (p Pool) GetJob(jobID int, userID int) (*Job, error) {
	for _, v := range p.jobs {
		if v.UserID == userID && v.ID == jobID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Job not found")
}

func (p Pool) GetPoolSize() int {
	return len(p.jobs)
}

type Job struct {
	ID         int           `json:"id"`
	Name       string        `json:"jobName"`
	URL        string        `json:"jobUrl"`
	Frequency  int64         `json:"frequency"`
	UserID     int           `json:"userID"`
	Active     bool          `json:"active"`
	CreatedAt  time.Time     `json:"createdAt"`
	ModifiedAt time.Time     `json:"modifiedAt"`
	status     string        `json:"-"`
	Done       chan bool     `json:"-"`
	logger     logger.Logger `json:"-"`
}

func NewJob(jobId int, userID int, name, url string, freq int64) *Job {

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
		Done:       make(chan bool),
	}

	return &j
}

func (j *Job) Run() {
	log.Printf("Starting job | ID: %d | Name: %s\n", j.ID, j.Name)
	j.SetActive(true)
	go j.run()
}

func (j *Job) run() {
	for {
		select {
		case <-j.Done:
			log.Printf("Ending job | ID: %d | Name: %s | URL: %s\n", j.ID, j.Name, j.URL)
			j.SetStatus("NA")
			j.SetActive(false)
			j.SetModifiedNow()
			return
		default:
			err := j.URLStatus()
			if err != nil {
				log.Println(err)
				j.SetStatus("DOWN")
			} else {
				j.SetStatus("UP")
			}
		}
	}
}

func (j *Job) SetModifiedNow() {
	j.ModifiedAt = time.Now()
}

func (j Job) URLStatus() error {
	resp, err := http.Get(j.URL)
	if err == nil && resp.StatusCode == 200 {
		return nil
	} else {
		return err
	}
}

func (j *Job) Stop() {
	j.Done <- true
}

func (j *Job) SetStatus(s string) {
	j.status = s
}

func (j *Job) SetActive(a bool) {
	j.Active = a
}

func (j *Job) GetStatus() string {
	return j.status
}

func (j *Job) GetActive() bool {
	return j.Active
}
