package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	modelUser "wmt/internal/model/user"
	"wmt/logger"
	serviceRoom "wmt/service/room"

	"github.com/gorilla/websocket"
)

// Manager 所有 websocket 信息
type Manager struct {
	Group                   map[string]map[string]*Client
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
}

// Client 单个 websocket 信息
type Client struct {
	ClientId, Group string
	User            modelUser.User
	Socket          *websocket.Conn
	Message         chan []byte
}

// messageData 单个发送数据信息
type MessageData struct {
	Id, Group string
	Message   []byte
}

// groupMessageData 组广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// 广播发送数据信息
type BroadCastMessageData struct {
	Message []byte
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Printf("client [%s] disconnect", *c.User.FirstName)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.ClientId, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		log.Printf("client [%s] receive message: %s", *c.User.FirstName, string(message))
		// 处理拿到的消息
		err = handleMessage(c, message)
		if err != nil {
			// log.Printf("client [%s] writemessage err: %s", c.ClientId, err)
			log.Printf("client [%s] writemessage err: %s", *c.User.FirstName, err)
		}
		// c.Message <- message
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		log.Printf("client [%s] disconnect", *c.User.FirstName)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", *c.User.FirstName, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte("websocker error"))
				return
			}
			log.Printf("client [%s] send message: %s ", *c.User.FirstName, string(message))
			c.Socket.WriteJSON(string(message))
		}
	}
}

// 启动 websocket 管理器
func (manager *Manager) Start() {
	logger.Debuglog.Printf("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			log.Printf("client [%s] connect", *client.User.FirstName)
			log.Printf("register client [%s] to group [%s]", client.User.FirstName, client.Group)
			manager.Lock.Lock()
			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*Client)
				manager.groupCount += 1
			}
			manager.Group[client.Group][client.ClientId] = client
			manager.clientCount += 1
			manager.Lock.Unlock()
			var enterMsg = struct {
				Type       string         `json:"type"`
				User       modelUser.User `json:"user"`
				GroupCount int            `json:"groupCount"`
			}{Type: "enter", User: client.User, GroupCount: len(manager.Group[client.Group])}
			enterMsgBytes, _ := json.Marshal(enterMsg)
			// 加入房间
			serviceRoom.EnterRoom(client.Group, client.User.UserId)
			manager.SendGroup(client.Group, enterMsgBytes)
			log.Printf("client count: %d", WebsocketManager.clientCount)
		// 注销
		case client := <-manager.UnRegister:
			log.Printf("unregister client [%s] from group [%s]", client.ClientId, client.Group)
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.ClientId]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.ClientId)
					manager.clientCount -= 1
					if len(manager.Group[client.Group]) == 0 {
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
				}
			}
			manager.Lock.Unlock()
			var leaveMsg = struct {
				Type       string         `json:"type"`
				User       modelUser.User `json:"user"`
				GroupCount int            `json:"groupCount"`
			}{Type: "leave", User: client.User, GroupCount: len(manager.Group[client.Group])}
			leaveMsgBytes, _ := json.Marshal(leaveMsg)
			manager.SendGroup(client.Group, leaveMsgBytes)
			log.Printf("client count: %d", WebsocketManager.clientCount)
			// 从房间移出用户
			serviceRoom.LeaveRoom(client.Group, client.User.UserId)
			// 房间无用户则解散房间
			// if len(manager.Group[client.Group]) == 0 {
			// 	serviceRoom.DeleteRoom(client.Group)
			// }
		}
	}
}

// 处理单个 client 发送数据
func (manager *Manager) SendService() {
	for data := range manager.Message {
		if groupMap, ok := manager.Group[data.Group]; ok {
			if conn, ok := groupMap[data.Id]; ok {
				conn.Message <- data.Message
			}
		}
	}
}

// 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for data := range manager.GroupMessage {
		if groupMap, ok := manager.Group[data.Group]; ok {
			for _, conn := range groupMap {
				conn.Message <- data.Message
			}
		}
	}
}

// 处理广播数据
func (manager *Manager) SendAllService() {
	for data := range manager.BroadCastMessage {
		for _, v := range manager.Group {
			for _, conn := range v {
				conn.Message <- data.Message
			}
		}
	}
}

// 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, message []byte) {
	logger.Debuglog.Printf("send message to client [%s]", id)
	data := &MessageData{
		Id:      id,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group:   group,
		Message: message,
	}
	manager.GroupMessage <- data
}

// 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message: message,
	}
	manager.BroadCastMessage <- data
}

// 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

// 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group:            make(map[string]map[string]*Client),
	Register:         make(chan *Client, 128),
	UnRegister:       make(chan *Client, 128),
	GroupMessage:     make(chan *GroupMessageData, 128),
	Message:          make(chan *MessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
	groupCount:       0,
	clientCount:      0,
}

// 处理收到的数据
func handleMessage(c *Client, message []byte) error {
	log.Printf("handleMessage %v", string(message))
	if groupMap, ok := WebsocketManager.Group[c.Group]; ok {
		for _, conn := range groupMap {
			if c.ClientId != conn.ClientId {
				conn.Message <- message
			}
		}
	}
	return nil
}

// 测试组广播
func TestSendGroup() {
	for {
		time.Sleep(time.Second * 2)
		WebsocketManager.SendGroup("abcd11", []byte("SendGroup message ----"+time.Now().Format("2006-01-02 15:04:05")))
	}
}

// 测试广播
func TestSendAll() {
	for {
		time.Sleep(time.Second * 5)
		WebsocketManager.SendAll([]byte("SendAll message ----" + time.Now().Format("2006-01-02 15:04:05")))
		fmt.Println(WebsocketManager.Info())
	}
}
