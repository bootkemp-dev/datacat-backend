package router

import (
	"fmt"

	"github.com/bootkemp-dev/datacat-backend/auth"
	"github.com/bootkemp-dev/datacat-backend/handlers"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	authRouter := r.Group("auth")
	{
		authRouter.POST("/register", handlers.Register)
		authRouter.POST("/login", handlers.Login)
	}
	auth2 := r.Group("protected")
	{
		auth2.Use(auth.AuthMiddleware())
		auth2.GET("/me", handlers.Me)
		auth2.POST("/refresh", handlers.Refresh)
		auth2.POST("/logout", handlers.Logout)

		//background jobs
		auth2.POST("/jobs", handlers.AddJob)
		auth2.GET("/jobs", handlers.GetJobs)
		auth2.GET("/job/:id/status", handlers.GetJobstatus)
		auth2.GET("/job/:id/active", handlers.GetJobActive)
		auth2.POST("/job/:id/pause", handlers.PauseJob)
		auth2.POST("/job/:id/restart", handlers.RestartJob)
		auth2.DELETE("/job/:id", handlers.DeleteJob)
	}

	return r
}

func Run(port string) {
	router := setupRouter()
	router.Run(fmt.Sprintf(":%s", port))
}
