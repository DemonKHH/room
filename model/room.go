package model

type Room struct {
	RoomId   string `json:"roomId"`
	RoomName string `json:"roomName"`
	VideoUrl string `json:"videoUrl"`
	Members  []User `json:"members"`
}
