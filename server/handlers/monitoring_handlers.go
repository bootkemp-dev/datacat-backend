package handlers

import (
	"net/http"
	"time"

	"github.com/bootkemp-dev/datacat-backend/database"
	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/bootkemp-dev/datacat-backend/monitoring"
	"github.com/bootkemp-dev/datacat-backend/utils"

	"github.com/gin-gonic/gin"
)

var jobPool monitoring.Pool

func init() {
	jobPool = monitoring.NewPool()
}

func AddJob(c *gin.Context) {
	var request models.NewJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := utils.ValidateNewJob(request.JobName, request.JobURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "id not set in context",
		})
		return
	}

	jobID, err := database.InsertNewJob(request.JobName, request.JobURL, request.Frequency, id.(float64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	j := monitoring.NewJob(jobID, id.(float64), request.JobName, request.JobURL, time.Duration(request.Frequency))
	jobPool.AddJob(j)
	j.Run()

	c.JSON(http.StatusOK, gin.H{
		"id":   jobID,
		"name": request.JobName,
		"url":  request.JobURL,
	})
	return
}

func GetAllJobs(c *gin.Context) {

}

func DeleteJob(c *gin.Context) {

}
