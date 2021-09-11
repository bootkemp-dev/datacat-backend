package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := getTokenFromHeader(c.Request)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		username, id, err := isTokenValid(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Not Authorized",
			})
			c.Abort()
			return
		} else {
			c.Set("id", id)
			c.Set("username", username)
			c.Next()
		}
	}
}

func getTokenFromHeader(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("Not Valid")
	}
	tokenString := strings.TrimSpace(splitToken[1])
	return tokenString, nil
}
