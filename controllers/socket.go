package controllers

import (
	"log"
	"net/http"

	modelUser "wmt/internal/model/user"
	ws "wmt/lib/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// gin 处理 websocket handler
func WsClient(ctx *gin.Context) {
	firstName := ctx.GetString("firstName")
	userId := ctx.GetString("userId")
	email := ctx.GetString("email")
	user := modelUser.User{
		FirstName: &firstName,
		UserId:    userId,
		Email:     &email,
	}
	group := ctx.Param("channel")
	log.Printf("group: %v", group)
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("websocket connect error: %s", err)
		return
	}
	client := &ws.Client{
		ClientId: uuid.New().String(),
		User:     user,
		Group:    group,
		Socket:   conn,
		Message:  make(chan []byte, 1024),
	}
	log.Printf("ct client = %v ", client)
	ws.WebsocketManager.RegisterClient(client)
	go client.Read()
	go client.Write()
}
