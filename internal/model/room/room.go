package modelRoom

import (
	modelUser "wmt/internal/model/user"
)

type Room struct {
	RoomId   string           `json:"roomId"`
	RoomName string           `json:"roomName"`
	VideoUrl string           `json:"videoUrl"`
	Members  []modelUser.User `json:"members"`
}
