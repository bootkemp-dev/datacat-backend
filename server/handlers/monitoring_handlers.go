package handlers

import (
	"database/sql"
	"log"
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
	// load existing jobs

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

	jobID, err := database.InsertNewJob(request.JobName, request.JobURL, request.Frequency, id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	j := monitoring.NewJob(jobID, id.(int), request.JobName, request.JobURL, time.Duration(request.Frequency))
	jobPool.Jobs = append(jobPool.Jobs, j)
	j.Run()

	c.JSON(http.StatusOK, gin.H{
		"id":   jobID,
		"name": request.JobName,
		"url":  request.JobURL,
	})
	return
}

func GetAllJobs(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "id not set in context",
		})
		return
	}

	jobs, err := database.GetAllJobs(userID.(int))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "No jobs found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//check active jobs in the pool and assign them their status
	for i := range jobs {
		if jobs[i].Active == true {
			pj, err := jobPool.GetJob(jobs[i].ID, userID.(int))
			if err != nil {
				log.Println(err)
				continue
			}
			jobs[i].Status = pj.GetStatus()
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
	return
}

func DeleteJob(c *gin.Context) {

}
