package handlers

import (
	"net/http"

	"github.com/bootkemp-dev/datacat-backend/database"
	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/bootkemp-dev/datacat-backend/utils"
	"github.com/gin-gonic/gin"
)

//Register takes new user as a request, validates the data and inserts it into the database
func Register(c *gin.Context) {
	var request models.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//validate the data
	err := utils.NewUserValidator(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//check if username and email are already in the database
	err = database.CheckIfUsernameExists(request.Username)
	if err != nil {
		if err.Error() == "username exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err = database.CheckIfEmailExists(request.Email)
	if err != nil {
		if err.Error() == "email exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//create password hash
	hashedPassword, err := utils.HashPassword(request.Password1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//insert user
	err = database.InsertUser(request.Username, request.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.Status(200)
	return
}

func Login(c *gin.Context) {

}

func Me(c *gin.Context) {

}

func Logout(c *gin.Context) {

}
