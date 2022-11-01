package ws

type WsMsg struct {
	ClientId string `json:"clientId"`
	Type     string `json:"type"`
}
