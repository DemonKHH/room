package api

import (
	"encoding/json"
	"log"
	"net/http"
	"room/model"
	db "room/service/mongo"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 获取房间
func GetRooms(context *gin.Context) {
	var rooms []model.Room
	client := db.GetMongoClient()
	rltBytes, err := db.Find(client, "rooms", bson.M{})
	if err != nil {
		log.Printf("查找房间发生错误: %v", err)
		context.JSON(http.StatusOK, model.CommonResponse{
			Code: 1,
			Msg:  err.Error(),
			Data: make([]model.Room, 0),
		})
		return
	}
	json.Unmarshal(rltBytes, &rooms)
	log.Printf("rltBytes %v", rooms)
	if len(rooms) == 0 {
		rooms = make([]model.Room, 0)
	}
	context.JSON(http.StatusOK, model.CommonResponse{
		Code: 0,
		Msg:  "get successfully",
		Data: rooms,
	})
}

// // 加入房间
// func AddRoom(context *gin.Context) {

// }

// // 离开房间
// func LeaveRoom(context *gin.Context) {
// 	// 如果房间无人则解散房间
// }
