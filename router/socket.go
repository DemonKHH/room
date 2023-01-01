package router

import (
	"wmt/controllers"

	"github.com/gin-gonic/gin"
)

func Socket(r *gin.Engine) {
	room := r.Group("/room")
	room.GET("/:channel", controllers.WsClient)
}
