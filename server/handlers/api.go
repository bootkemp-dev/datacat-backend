package handlers

import (
	"log"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/bootkemp-dev/datacat-backend/database"
	"github.com/bootkemp-dev/datacat-backend/logger"
	"github.com/bootkemp-dev/datacat-backend/mailing"
	"github.com/bootkemp-dev/datacat-backend/models"
)

type API struct {
	database *database.Database
	jobPool  models.Pool
	mailing  *mailing.Mailing
	logger   *logger.Logger
}

func NewApi(c config.Config) (*API, error) {
	db, err := database.NewDatabase(c)
	if err != nil {
		return nil, err
	}

	m, err := mailing.NewMailing(c)
	if err != nil {
		//switch later to return nil, err
	}

	l, err := logger.NewLogger(&c)
	if err != nil {
		return nil, err
	}

	defer l.Close()

	jobPool := models.NewPool()

	api := API{
		database: db,
		jobPool:  jobPool,
		mailing:  m,
		logger:   l,
	}
	go api.populateJobPool()

	return &api, nil
}

func (a *API) populateJobPool() error {
	log.Println("Populating job pool...")
	jobs, err := a.database.GetAllJobs()
	if err != nil {
		log.Printf("Populating job pool failed: %v\n", err)
		return err
	}
	log.Printf("Found %d jobs\n", len(jobs))

	for _, j := range jobs {
		a.jobPool.AddJob(j)
		if j.Active {
			j.Run()
		}
	}

	return nil
}
