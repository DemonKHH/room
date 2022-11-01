package model

type MediaInfo struct {
	Type          string `json:"type"`
	MediaUrl      string `json:"mediaUrl"`
	MediaStatus   string `json:"mediaStatus"`
	MediaDuration string `json:"mediaDuration"`
}
