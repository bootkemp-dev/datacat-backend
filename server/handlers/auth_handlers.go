package handlers

import (
	"database/sql"
	"net/http"

	"github.com/bootkemp-dev/datacat-backend/auth"
	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/bootkemp-dev/datacat-backend/utils"
	"github.com/gin-gonic/gin"
)

//Register takes new user as a request, validates the data and inserts it into the database
func (a *API) Register(c *gin.Context) {
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
	err = a.database.CheckIfUsernameExists(request.Username)
	if err != sql.ErrNoRows {
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Username already exists in the database",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err = a.database.CheckIfEmailExists(request.Email)
	if err != sql.ErrNoRows {
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Email already exists in the database",
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
	err = a.database.InsertUser(request.Username, request.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.Status(201)
	return
}

func (a *API) Login(c *gin.Context) {
	var request models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//check if the username is in the database
	err := a.database.CheckIfUsernameExists(request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
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

	//get id and password
	passwordHash, id, err := a.database.GetIDAndPasswordHash(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//compare password and hash
	err = utils.CompareHashAndPassword(passwordHash, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	//generate token
	token, _, err := auth.GenerateToken(request.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	/*
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  *exp,
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
		})
	*/
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
	return
}

func (a *API) Me(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "username not set in context",
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

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"username": username,
	})
	return
}

func (a *API) Refresh(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "username not set in context",
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

	token, exp, err := auth.GenerateToken(username.(string), id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Expires:  *exp,
		Value:    token,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})
}

func (a *API) Logout(c *gin.Context) {

}
