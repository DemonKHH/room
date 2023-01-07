package router

import (
	"wmt/controllers"

	"github.com/gin-gonic/gin"
)

func Room(r *gin.Engine) {
	r.GET("/getRooms", controllers.GetRooms)
	r.GET("/getRoomInfoByRoomId", controllers.GetRoomInfoByRoomId)
	r.POST("/createRoom", controllers.CreateRoom)
	r.POST("/deleteRoom", controllers.DeleteRoom)
	r.POST("/enterRoom", controllers.EnterRoom)
	r.POST("/leaveRoom", controllers.LeaveRoom)
}
