package router

import (
	"wmt/controllers"

	"github.com/gin-gonic/gin"
)

func User(r *gin.Engine) {
	r.POST("/signup", controllers.Signup())
	r.POST("/login", controllers.Login())
}
