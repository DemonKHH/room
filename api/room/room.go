package api

import (
	"log"
	"net/http"
	"room/model"

	serviceRoom "room/service/room"

	"github.com/gin-gonic/gin"
)

// 获取房间
func GetRooms(context *gin.Context) {
	rooms, err := serviceRoom.GetRooms()
	if err != nil {
		log.Printf("查找房间发生错误: %v", err)
		context.JSON(http.StatusOK, model.CommonResponse{
			Code: 1,
			Msg:  err.Error(),
			Data: make([]model.Room, 0),
		})
		return
	}
	log.Printf("rooms: %v", rooms)
	context.JSON(http.StatusOK, model.CommonResponse{
		Code: 0,
		Msg:  "get successfully",
		Data: rooms,
	})
}

func GetRoomInfoByRoomId(context *gin.Context) {
	roomId := context.Query("roomId")
	rooms, err := serviceRoom.GetRoomInfoByRoomId(roomId)
	if err != nil {
		log.Printf("查找房间发生错误: %v", err)
		context.JSON(http.StatusOK, model.CommonResponse{
			Code: 1,
			Msg:  err.Error(),
			Data: make([]model.Room, 0),
		})
		return
	}
	if len(rooms) == 0 {
		context.JSON(http.StatusOK, model.CommonResponse{
			Code: 0,
			Msg:  "未找到此房间",
			Data: rooms,
		})
		return
	}
	context.JSON(http.StatusOK, model.CommonResponse{
		Code: 0,
		Msg:  "成功获取房间信息",
		Data: rooms,
	})
}

// 进入房间
func EnterRoom(context *gin.Context) {

}

// 创建房间
func CreateRoom(context *gin.Context) {
	var room model.Room
	err := context.ShouldBindJSON(&room)
	log.Printf("room body: %v", room)
	if err != nil {
		context.JSON(http.StatusOK, model.CommonResponse{
			Code: 1,
			Msg:  "请求参数错误",
			Data: err.Error(),
		})
	}
	err = serviceRoom.CreateRoom(room.RoomId, room.RoomName, room.VideoUrl)
	if err != nil {
		context.JSON(http.StatusOK, model.CommonResponse{
			Code: 1,
			Msg:  "创建房间失败",
			Data: err.Error(),
		})
	}
	context.JSON(http.StatusOK, model.CommonResponse{
		Code: 0,
		Msg:  "创建房间成功",
	})
}

// // 离开房间
// func LeaveRoom(context *gin.Context) {
// 	// 如果房间无人则解散房间
// }
