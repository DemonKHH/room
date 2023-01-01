package router

import (
	"log"
	"wmt/middleware"

	ws "wmt/lib/websocket"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.Engine) {
	log.Printf("[router] load routes")
	User(r)
	Auth(r)
	Room(r)
	Socket(r)
}

func InitRoutes() {
	router := gin.Default()
	router.Use(middleware.Cors())
	LoadRoutes(router)
	// router.StaticFS("/web", http.Dir("web"))
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendGroupService()
	go ws.WebsocketManager.SendAllService()
	router.Run(":8999")
}
