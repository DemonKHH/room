package controllers

import (
	"log"
	"net/http"
	modelRoom "wmt/internal/model/room"
	response "wmt/pkg/common/response"

	serviceRoom "wmt/service/room"

	"github.com/gin-gonic/gin"
)

// 获取房间
func GetRooms(context *gin.Context) {
	rooms, err := serviceRoom.GetRooms()
	if err != nil {
		log.Printf("查找房间发生错误: %v", err)
		context.JSON(http.StatusOK, response.ResponseMsg{
			Code: 1,
			Msg:  err.Error(),
			Data: []string{},
		})
		return
	}
	log.Printf("rooms: %v", rooms)
	context.JSON(http.StatusOK, response.ResponseMsg{
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
		context.JSON(http.StatusOK, response.ResponseMsg{
			Code: 1,
			Msg:  err.Error(),
			Data: []string{},
		})
		return
	}
	context.JSON(http.StatusOK, response.ResponseMsg{
		Code: 0,
		Msg:  "成功获取房间信息",
		Data: rooms,
	})
}

// 进入房间
func EnterRoom(context *gin.Context) {
	var room modelRoom.Room
	err := context.ShouldBindJSON(&room)
	if err != nil {
		context.JSON(http.StatusOK, response.ResponseMsg{
			Code: 1,
			Msg:  "请求参数错误",
			Data: err.Error(),
		})
		return
	}
	userId := context.GetString("userId")
	log.Printf("enterRoom roomId %v \n userId %v\n", room.RoomId, userId)
	err = serviceRoom.EnterRoom(room.RoomId, userId)
	if err != nil {
		log.Printf("进入房间发生错误: %v", err)
		context.JSON(http.StatusOK, response.ResponseMsg{
			Code: 1,
			Msg:  err.Error(),
			Data: []string{},
		})
		return
	}
	context.JSON(http.StatusOK, response.ResponseMsg{
		Code: 0,
		Msg:  "成功进入房间",
	})
}

// 创建房间
func CreateRoom(context *gin.Context) {
	var room modelRoom.Room
	err := context.ShouldBindJSON(&room)
	log.Printf("room body: %v", room)
	if err != nil {
		context.JSON(http.StatusOK, response.ResponseMsg{
			Code: 1,
			Msg:  "请求参数错误",
			Data: err.Error(),
		})
		return
	}
	err = serviceRoom.CreateRoom(room.RoomId, room.RoomName, room.VideoUrl)
	if err != nil {
		context.JSON(http.StatusOK, response.ResponseMsg{
			Code: 1,
			Msg:  "创建房间失败",
			Data: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.ResponseMsg{
		Code: 0,
		Msg:  "创建房间成功",
	})
}

// // 离开房间
func LeaveRoom(context *gin.Context) {
	roomId := context.Query("roomId")
	userId := context.GetString("userId")
	log.Printf("leaveRoom roomId %v \n userId %v\n", roomId, userId)
	err := serviceRoom.LeaveRoom(roomId, userId)
	if err != nil {
		log.Printf("离开房间发生错误: %v", err)
		context.JSON(http.StatusOK, response.ResponseMsg{
			Code: 1,
			Msg:  err.Error(),
			Data: []string{},
		})
		return
	}
	context.JSON(http.StatusOK, response.ResponseMsg{
		Code: 0,
		Msg:  "已离开房间",
	})
}
