package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

// Client represents a client connected to a websocket server.
// It contains three fields: id, socket, and send.
// Client 代表连接到 websocket 服务器的客户端。
// 它包含三个字段：id、socket 和 send。
type Client struct {
	id     string          // 唯一的 ID
	socket *websocket.Conn // websocket 连接
	send   chan []byte     // 用于发送消息的通道，它是一个字节切片，用于存储待发送给这个客户端的信息
}

// read reads messages from the websocket connection of the client.
// It handles errors and broadcasts received messages to other clients.
// When an error occurs or the connection is closed, it closes the socket and unregisters the client from the manager.
// read 从客户端的 websocket 连接中读取消息。
// 它处理错误并将接收到的消息广播给其他客户端。
// 发生错误或关闭连接时，它将关闭套接字并从管理器中注销客户端。
func (c *Client) read() {
	defer func() { // 关闭套接字并从管理器中注销客户端
		manager.unregister <- c
		_ = c.socket.Close()
	}()

	for { // 无限循环，从客户端的 websocket 连接中读取消息
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c // 如果发生错误，则关闭套接字并从管理器中注销客户端
			_ = c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)}) // 将消息转换为 JSON 格式
		manager.broadcast <- jsonMessage
	}
}

// write writes messages to the websocket connection of the client.
// It handles closing the socket and sending messages.
// If the channel `c.send` is closed, it will send a close message to the websocket connection.
// For each message received from the channel, it sends a text message over the websocket.
// write 将消息写入客户端的 websocket 连接。
// 它处理关闭套接字和发送消息。
// 如果通道 c.send 被关闭，它将向 websocket 连接发送一个关闭消息。
// 对于从通道接收到的每条消息，它都会通过 websocket 发送一条文本消息。
func (c *Client) write() {
	defer func() { // 关闭套接字
		_ = c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send: // 从通道接收到的每条消息
			if !ok { // 如果通道 c.send 被关闭，它将向 websocket 连接发送一个关闭消息。
				c.socket.WriteMessage(websocket.CloseMessage, []byte{}) // 如果通道 c.send 被关闭，它将向 websocket 连接发送一个关闭消息。
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message) // 对于从通道接收到的每条消息，它都会通过 websocket 发送一条文本消息。
		}
	}
}
