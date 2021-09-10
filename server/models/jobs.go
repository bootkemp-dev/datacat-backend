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
	Name       string        `json:"job_name"`
	URL        string        `json:"job_url"`
	Frequency  int64         `json:"frequency"`
	UserID     int           `json:"user_id"`
	Active     bool          `json:"active"`
	CreatedAt  time.Time     `json:"created_at"`
	ModifiedAt time.Time     `json:"modified_at"`
	status     string        `json:"-"`
	done       chan struct{} `json:"-"`
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
		done:       make(chan struct{}),
	}

	return &j
}

func (j *Job) Run() {
	log.Printf("Starting job | ID: %d | Name: %s\n", j.ID, j.Name)
	j.SetActive(true)
	go func() {
		for {
			select {
			case <-j.done:
				log.Printf("Ending job | ID: %d | Name: %s | URL: %s\n", j.ID, j.Name, j.URL)
				j.SetStatus("NA")
				j.SetActive(false)
				return
			default:
				log.Printf("Job with ID: %d checking status of: %s ", j.ID, j.URL)
				err := j.URLStatus()
				if err != nil {
					log.Println(err)
					j.SetStatus("DOWN")
				} else {
					j.SetStatus("UP")
				}
			}
			time.Sleep(time.Duration(j.Frequency))
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

func (j Job) Stop() {
	j.done <- struct{}{}
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
