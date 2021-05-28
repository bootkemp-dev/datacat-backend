package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Not authenticated",
			})
			c.Abort()
			return
		}

		tokenString := cookie.Value
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
