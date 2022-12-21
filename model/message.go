package model

type BaseMessage struct {
	Type string          `json:"type"`
	Data BaseMessageData `json:"data"`
}

type BaseMessageData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type MediaInfo struct {
	Type          string `json:"type"`
	MediaUrl      string `json:"mediaUrl"`
	MediaStatus   string `json:"mediaStatus"`
	MediaDuration string `json:"mediaDuration"`
}
