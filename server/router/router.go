package router

import (
	"fmt"
	"log"

	"github.com/bootkemp-dev/datacat-backend/auth"
	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/bootkemp-dev/datacat-backend/handlers"
	"github.com/gin-gonic/gin"
)

func setupRouter(c config.Config) *gin.Engine {
	api, err := handlers.NewApi(c)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	authRouter := r.Group("auth")
	{
		authRouter.POST("/register", api.Register)
		authRouter.POST("/login", api.Login)
	}
	auth2 := r.Group("protected")
	{
		auth2.Use(auth.AuthMiddleware())
		auth2.GET("/me", api.Me)
		auth2.POST("/refresh", api.Refresh)
		auth2.POST("/logout", api.Logout)

		//background jobs
		auth2.POST("/jobs", api.AddJob)
		auth2.GET("/jobs", api.GetJobsFromPool)
		auth2.GET("/job/:id/status", api.GetJobstatus)
		auth2.GET("/job/:id/active", api.GetJobActive)
		auth2.POST("/job/:id/pause", api.PauseJob)
		auth2.POST("/job/:id/restart", api.RestartJob)
		auth2.DELETE("/job/:id", api.DeleteJob)
		auth2.GET("job/:id/socket.io", api.SocketIOHandler)
		auth2.POST("job/:id/socket.io", api.SocketIOHandler)
		auth2.GET("/job/:id/logs", api.JobLogHandler)
	}

	accountsRouter := r.Group("accounts")
	{
		accountsRouter.POST("/reset_password", api.HandleResetPassword)
		accountsRouter.GET("/reset_password", api.HandleResetTokenValidation)
		accountsRouter.PUT("/update_password", api.HandlePasswordChangeAfterReset)
	}

	return r
}

func Run(c config.Config) {
	router := setupRouter(c)
	router.Run(fmt.Sprintf(":%s", c.Server.Port))
}
