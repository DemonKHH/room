package router

import (
	"wmt/controllers"
	"wmt/middleware"

	"github.com/gin-gonic/gin"
)

func Auth(r *gin.Engine) {
	r.Use(middleware.Authenticate())
	r.GET("/user/refreshToken", controllers.RefreshToken())
	r.GET("/getUsers", controllers.GetUsers())
	r.GET("/getUser", controllers.GetUser())
	// 需要鉴权的 api
	// auth := r.Group("/auth")
}
