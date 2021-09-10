package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/bootkemp-dev/datacat-backend/utils"
	"github.com/gorilla/websocket"

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

	j := models.NewJob(jobID, id.(int), request.JobName, request.JobURL, request.Frequency)
	a.jobPool.AddJob(j)
	j.Run()
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

	if job.GetActive() == true {
		job.Stop()
		job.Run()
	} else {
		job.Run()
	}
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

func (a *API) JobInfoWebsocket(c *gin.Context) {

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
	var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer ws.Close()

	for {

		t, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		//prepare message
		msgToSend := []byte(job.GetStatus())
		ws.WriteMessage(t, msgToSend)

		time.Sleep(time.Duration(job.Frequency))
	}
}
