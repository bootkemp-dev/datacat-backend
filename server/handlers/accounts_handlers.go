package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/bootkemp-dev/datacat-backend/mailing"
	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/bootkemp-dev/datacat-backend/utils"
	"github.com/gin-gonic/gin"
)

func (a *API) HandleResetPassword(c *gin.Context) {
	config, err := config.NewConfig("./config.yml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "username not set",
		})
		return
	}
	//check if username exists in the database and get email
	email, err := a.database.GetUserEmail(username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User does not exist",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	}

	//generate one time token with expiration date and insert it into the database
	token, err := utils.GenerateRandomToken(config.Accounts.ResetPasswordTokenLength)
	timeToAdd := config.Accounts.ResetPasswordTokenExpiration
	err = a.database.UpdateResetPasswordToken(username, token, time.Now().Local().Add(time.Duration(timeToAdd)))

	go mailing.SendResetPasswordEmail(username, email, token)
	c.Status(http.StatusOK)
}

func (a *API) HandlePasswordChangeAfterReset(c *gin.Context) {
	username := c.Query("usename")
	token := c.Query("token")

	if username == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "username and token not set as query params",
		})
		return
	}

	var request models.PasswordResetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if request.Password1 != request.Password2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "passwords do not match",
		})
		return
	}

	//get token from the database
	exp, err := a.database.GetResetPasswordTokenExpiration(username, token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusUnauthorized)
		} else {
			c.Status(http.StatusInternalServerError)
		}
	}

	if !utils.InTimeSpan(*exp) {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": "password reset token already expired",
		})
		return
	}

	//hash the password
	hashedPassword, err := utils.HashPassword(request.Password1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//update password in the database
	err = a.database.UpdatePasswordHash(username, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (a *API) HandleResetTokenValidation(c *gin.Context) {

	username := c.Query("usename")
	token := c.Query("token")

	if username == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "username and token not set as query params",
		})
		return
	}

	//get token from the database
	exp, err := a.database.GetResetPasswordTokenExpiration(username, token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(http.StatusUnauthorized)
		} else {
			c.Status(http.StatusInternalServerError)
		}
	}

	if !utils.InTimeSpan(*exp) {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"message": "password reset token already expired",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "password reset token is valid",
	})

}
