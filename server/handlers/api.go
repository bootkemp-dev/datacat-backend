package handlers

import (
	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/bootkemp-dev/datacat-backend/database"
	"github.com/bootkemp-dev/datacat-backend/models"
)

type API struct {
	database *database.Database
	jobPool  models.Pool
}

func NewApi(c config.Config) (*API, error) {
	db, err := database.NewDatabase(c)
	if err != nil {
		return nil, err
	}

	jobPool := models.NewPool()

	api := API{
		database: db,
		jobPool:  jobPool,
	}

	return &api, nil
}
