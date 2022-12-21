package serviceRoom

import (
	"encoding/json"
	"log"
	"room/model"
	db "room/service/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteRoom(roomId string) {
	client := db.GetMongoClient()
	err := db.Delete(client, "rooms", bson.M{
		"roomid": roomId,
	})
	if err != nil {
		print("创建房间失败 %v", err)
	}
}

func CreateRoom(roomId string) {
	client := db.GetMongoClient()
	err := db.Insert(client, "rooms", model.Room{
		RoomId:   roomId,
		RoomName: roomId + "测试",
		Members:  []model.User{},
	})
	if err != nil {
		print("创建房间失败 %v", err)
	}
}

func AddRoom(roomId string, clientId string) {
	var rooms = []model.Room{}
	client := db.GetMongoClient()
	filter := bson.D{{Key: "roomid", Value: roomId}}
	respBytes, err := db.Find(client, "rooms", filter)
	if err != nil {
		log.Printf("leave room error: %v", err)
		return
	}
	json.Unmarshal(respBytes, &rooms)
	if len(rooms) == 0 {
		// 没找到对应的房间,先创建房间再加入房间
		CreateRoom(roomId)
	}
	update := bson.M{"$push": bson.M{"members": bson.M{
		"clientId": clientId,
		"userName": clientId + "test",
		"avator":   "https://www.dmoe.cc/random.php",
	}}}
	err = db.Update(client, "rooms", filter, update)
	if err != nil {
		log.Printf("error updating %v", err)
	}
}

func LeaveRoom(roomId string, clientId string) {
	var rooms = []model.Room{}
	client := db.GetMongoClient()
	filter := bson.D{{Key: "roomid", Value: roomId}}
	respBytes, err := db.Find(client, "rooms", filter)
	if err != nil {
		log.Printf("leave room error: %v", err)
		return
	}
	json.Unmarshal(respBytes, &rooms)
	if len(rooms) == 0 {
		// 没找到对应的房间 不处理
		return
	} else {
		// 移出 members 中对应的用户
		update := bson.M{"$pull": bson.M{"members": bson.M{"clientId": clientId}}}
		err = db.Update(client, "rooms", filter, update)
		if err != nil {
			log.Printf("error updating %v", err)
		}
	}
}
