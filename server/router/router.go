package router

import (
	"fmt"

	"github.com/bootkemp-dev/datacat-backend/handlers"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	auth := r.Group("auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/me", handlers.Me)
	}
	return r
}

func Run(port string) {
	router := setupRouter()
	router.Run(fmt.Sprintf(":%s", port))
}
