package router

import (
	"log"
	"net/http"
	ws "room/lib/websocket"
	"room/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.Engine) {
	log.Printf("[router] load routes")
}

func InitRoutes() {
	router := gin.Default()
	router.Use(middleware.Cors())
	router.StaticFS("/web", http.Dir("web"))
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendGroupService()
	wsGroup := router.Group("/room")
	{
		wsGroup.GET("/:channel", ws.WebsocketManager.WsClient)
	}
	router.Run(":8090")
}
