package handlers

import (
	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/gin-gonic/gin"
)

func AddJob(c *gin.Context) {
	var request models.NewJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {

	}

}

func GetAllJobs(c *gin.Context) {

}

func DeleteJob(c *gin.Context) {

}
