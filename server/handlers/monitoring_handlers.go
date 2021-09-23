package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/bootkemp-dev/datacat-backend/utils"
	socketio "github.com/googollee/go-socket.io"

	"github.com/gin-gonic/gin"
)

func (a *API) AddJob(c *gin.Context) {
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

	jobID, err := a.database.InsertNewJob(request.JobName, request.JobURL, request.Frequency, id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	j, err := models.NewJob(jobID, id.(int), request.JobName, request.JobURL, request.Frequency, time.Now(), time.Now(), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	a.jobPool.AddJob(j)
	j.Run()

	go func() {
		logMessage := fmt.Sprintf("Job with ID %d has been created", jobID)
		if err := a.logger.WriteLogToFile(logMessage); err != nil {
			log.Println(err)
		}
		if err := a.database.InsertJobLog(id.(int), j.ID, j.GetStatus(), logMessage); err != nil {
			log.Println(err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"id":   jobID,
		"name": request.JobName,
		"url":  request.JobURL,
	})
	return
}

func (a *API) GetJobstatus(c *gin.Context) {
	id := c.Param("id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	job, err := a.jobPool.GetJob(jobID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  job.GetStatus(),
	})

	return
}

func (a *API) GetJobs(c *gin.Context) {

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	jobIDString := c.Query("id")
	if jobIDString != "" {
		jobID, err := strconv.Atoi(jobIDString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		job, err := a.database.GetJobByID(jobID, userID.(int))
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Job nor found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"job":     &job,
		})
		return
	} else {
		jobs, err := a.database.GetAllJobsByUserID(userID.(int))
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Job not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"job":     &jobs,
		})
		return
	}
}

func (a *API) PauseJob(c *gin.Context) {
	id := c.Param("id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	job, err := a.jobPool.GetJob(jobID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if job.GetActive() == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Can not pause a job which is not active",
		})
		return
	}

	go job.Stop()
	job.SetModifiedNow()

	go func() {
		logMessage := fmt.Sprintf("Job with ID: %d has been paused", job.ID)
		if err := a.logger.WriteLogToFile(logMessage); err != nil {
			log.Println(err)
		}

		if err := a.database.InsertJobLog(userID.(int), jobID, job.Status, logMessage); err != nil {
			log.Println(err)
		}
	}()

	err = a.database.UpdateJobActive(false, jobID, job.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (a *API) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	//delete job from the pool
	err = a.jobPool.RemoveJob(jobID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// delete job from the database
	err = a.database.DeleteJob(jobID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
	return
}

func (a *API) RestartJob(c *gin.Context) {
	id := c.Param("id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	//get job
	job, err := a.jobPool.GetJob(jobID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if job.GetActive() {
		job.Stop()
		job.Run()
		job.SetModifiedNow()
	} else {
		err := a.database.UpdateJobActive(true, jobID, userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		job.SetModifiedNow()
		job.Run()
	}

	go func() {
		logMessage := fmt.Sprintf("Job with ID: %d has been restarted", job.ID)
		if err := a.database.InsertJobLog(userID.(int), jobID, job.GetStatus(), logMessage); err != nil {
			log.Println(err)
		}
		if err := a.logger.WriteLogToFile(logMessage); err != nil {
			log.Println(err)
		}
	}()

	c.Status(200)
}

func (a *API) GetJobActive(c *gin.Context) {
	id := c.Param("id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	//get job
	job, err := a.jobPool.GetJob(jobID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"active":  job.GetActive(),
	})
}

func (a *API) JobLogHandler(c *gin.Context) {

	id := c.Param("id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var limit int
	var offset int

	limitQuery := c.Query("limit")
	if limitQuery == "" {
		limit = 0
	} else {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	offsetQuery := c.Query("offset")
	if offsetQuery == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	var logs []*models.JobLog

	if limit == 0 && offset == 0 {
		logs, err = a.database.GetAllJobLogsByID(jobID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"mesage": err.Error(),
			})
			return
		}

	} else {
		logs, err = a.database.GetJobLogsByID(jobID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"mesage": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"offset": offset,
		"limit":  limit,
	})
	return
}

func (a *API) GetJobsFromPool(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	jobs := a.jobPool.GetJobsByUserID(userID.(int))
	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})

	return
}

func (a *API) SocketIOHandler(c *gin.Context) {
	id := c.Param("id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "id not set in context",
		})
		return
	}

	//get job
	job, err := a.jobPool.GetJob(jobID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if !job.GetActive() {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Job not running",
		})
		return
	}

	a.SocketServer.OnEvent("/", "ping", func(s socketio.Conn) *time.Duration {
		return job.GetPing()
	})

	a.SocketServer.ServeHTTP(c.Writer, c.Request)
}
