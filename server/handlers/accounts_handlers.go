package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/bootkemp-dev/datacat-backend/database"
	"github.com/bootkemp-dev/datacat-backend/mailing"
	"github.com/bootkemp-dev/datacat-backend/utils"
	"github.com/gin-gonic/gin"
)

func HandleResetPassword(c *gin.Context) {
	config, err := config.NewConfig("./config.yml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	username := c.Query("username")
	//check if username exists in the database and get email
	email, err := database.GetUserEmail(username)
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
	err = database.UpdateResetPasswordToken(username, token, time.Now().Local().Add(time.Duration(timeToAdd)))

	go mailing.SendResetPasswordEmail(username, email, token)
	c.Status(http.StatusOK)
}

func HandlePasswordChangeAfterReset(c *gin.Context) {
	username := c.Query("usename")
	token := c.Query("token")

	if username == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "username and token not set as query params",
		})
		return
	}
}

func HandleResetTokenValidation(c *gin.Context) {

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
	exp, err := database.GetResetPasswordTokenExpiration(username, token)
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
