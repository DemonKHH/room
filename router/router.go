package router

import (
	"log"
	"net/http"
	api "room/api/room"
	ws "room/lib/websocket"
	"room/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.Engine) {
	log.Printf("[router] load routes")
	r.GET("/getRooms", api.GetRooms)
}

func InitRoutes() {
	router := gin.Default()
	LoadRoutes(router)
	router.Use(middleware.Cors())
	router.StaticFS("/web", http.Dir("web"))
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendGroupService()
	wsGroup := router.Group("/room")
	{
		wsGroup.GET("/:channel", ws.WebsocketManager.WsClient)
	}
	router.Run(":8999")
}
