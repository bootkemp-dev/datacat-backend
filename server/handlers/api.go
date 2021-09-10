package handlers

import (
	"log"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/bootkemp-dev/datacat-backend/database"
	"github.com/bootkemp-dev/datacat-backend/mailing"
	"github.com/bootkemp-dev/datacat-backend/models"
)

type API struct {
	database *database.Database
	jobPool  models.Pool
	mailing  *mailing.Mailing
}

func NewApi(c config.Config) (*API, error) {
	db, err := database.NewDatabase(c)
	if err != nil {
		return nil, err
	}

	m, err := mailing.NewMailing(c)
	if err != nil {
		log.Println(err)
		//switch later to return nil, err
	}

	jobPool := models.NewPool()

	api := API{
		database: db,
		jobPool:  jobPool,
		mailing:  m,
	}

	return &api, nil
}
