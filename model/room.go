package model

type Room struct {
	RoomId   string `json:"roomId"`
	RoomName string `json:"roomName"`
	Members  []User `json:"members"`
}
