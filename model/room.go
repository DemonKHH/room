package model

type Room struct {
	RoomId      int    `json:"roomId"`
	RoomName    string `json:"roomName"`
	RoomMembers []User `json:"roomMembers"`
}
