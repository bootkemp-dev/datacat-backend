package router

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	r := gin.Default()
	auth := r.Group()
	{
		auth.POST()
	}
	return r
}

func Run() {

}
