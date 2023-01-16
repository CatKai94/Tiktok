package router

import (
	"github.com/gin-gonic/gin"
	"tiktokbackend/controller/user"
	"tiktokbackend/service"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	apiRouter.GET("/saveuser", service.SaveUser)
	apiRouter.GET("/saveuser", user.TestController)
	apiRouter.GET("/getuser", service.GetUser)
	apiRouter.GET("/updateuser", service.UpdateUser)
	apiRouter.GET("/deleteuser", service.DeleteUser)
}
