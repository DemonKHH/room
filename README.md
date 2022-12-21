# room
chat room base on websocket

1. 字段设计
 1）用户信息
 2）通信消息


统一通信信息
{
    type: "数据类型", // 例如 message
    data: {
        type: "数据类型", // 例如mediaInfo、userInfo
        data: "自定义需要传输的数据类型"
    }
}