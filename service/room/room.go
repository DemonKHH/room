package serviceRoom

import (
	"context"
	"encoding/json"
	"log"
	"time"
	modelRoom "wmt/internal/model/room"
	db "wmt/service/mongo"
	serviceUser "wmt/service/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var roomCollection *mongo.Collection = db.OpenCollection(db.GetMongoClient(), "rooms")

func GetRooms() (rooms []modelRoom.Room, err error) {
	client := db.GetMongoClient()
	rltBytes, err := db.Find(client, "rooms", bson.M{})
	if err != nil {
		log.Printf("查找房间发生错误: %v", err)
	}
	json.Unmarshal(rltBytes, &rooms)
	log.Printf("rltBytes %v", rooms)
	if len(rooms) == 0 {
		rooms = make([]modelRoom.Room, 0)
	}
	return rooms, err
}

func GetRoomInfoByRoomId(roomId string) (room modelRoom.Room, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = roomCollection.FindOne(ctx, bson.M{"roomid": roomId}).Decode(&room)
	return room, err
}

func DeleteRoom(roomId string) {
	client := db.GetMongoClient()
	err := db.Delete(client, "rooms", bson.M{
		"roomid": roomId,
	})
	if err != nil {
		print("创建房间失败 %v", err)
	}
}

func CreateRoom(roomId string, roomName string, videoUrl string) error {
	var err error
	client := db.GetMongoClient()
	err = db.Insert(client, "rooms", bson.M{
		"roomid":   roomId,
		"roomname": roomName,
		"videourl": videoUrl,
		"members":  []string{},
	})
	if err != nil {
		print("创建房间失败 %v", err)
	}
	return err
}

func EnterRoom(roomId string, userId string) error {
	client := db.GetMongoClient()
	filter := bson.D{{Key: "roomid", Value: roomId}}
	user, err := serviceUser.GetUser(userId)
	if err != nil {
		log.Printf("get user info error: %v", err)
		return err
	}
	_, err = GetRoomInfoByRoomId(roomId)
	if err != nil {
		log.Printf("get room info error: %v", err)
		return err
	}
	update := bson.M{"$addToSet": bson.M{"members": bson.M{
		"userid":   user.UserId,
		"username": user.FirstName,
		"avator":   user.Avator,
	}}}
	err = db.Update(client, "rooms", filter, update)
	if err != nil {
		log.Printf("error updating %v", err)
	}
	return err
}

func LeaveRoom(roomId string, userId string) error {
	client := db.GetMongoClient()
	filter := bson.D{{Key: "roomid", Value: roomId}}

	// 移出 members 中对应的用户
	update := bson.M{"$pull": bson.M{"members": bson.M{"userid": userId}}}
	err := db.Update(client, "rooms", filter, update)
	if err != nil {
		log.Printf("error updating %v", err)
	}
	return err
}
