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
	auth2 := r.Group("auth2")
	{
		auth2.Use(auth.AuthMiddleware())
		auth2.GET("/me", handlers.Me)
		auth2.POST("/logout", handlers.Logout)
	}

	return r
}

func Run(port string) {
	router := setupRouter()
	router.Run(fmt.Sprintf(":%s", port))
}
