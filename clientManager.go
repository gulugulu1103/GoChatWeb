package main

import "encoding/json"

// ClientManager manages client connections and message broadcasting.
// It keeps track of connected clients, handles new client registration,
// and sends received messages to all connected clients.

type ClientManager struct {
	clients    map[*Client]bool // 用于存储连接的客户端的映射
	broadcast  chan []byte      // 用于广播消息的通道
	register   chan *Client     // 用于注册客户端的通道
	unregister chan *Client     // 用于注销客户端的通道
}

// send is a method of the ClientManager struct. It sends the provided message to all clients in the 'clients' map except for the 'ignore' client.
//
// Parameters:
// - message: The message to be sent, as a byte slice.
// - ignore: The client to be ignored, as a pointer to a Client struct.
// send 是 ClientManager 结构体的方法。它将提供的消息发送到 clients 映射中除 ignore 客户端之外的所有客户端。
// 参数：
// - message：要发送的消息，作为字节切片。
// - ignore：要忽略的客户端，作为指向 Client 结构体的指针。
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

// start is a method of the ClientManager struct. It starts a for loop that continuously listens for incoming events on multiple channels.
// When a connection is received on the register channel, the connection is added to the clients map and a JSON message indicating a new socket connection is sent to the client.
// When a connection is received on the unregister channel, the connection is checked if it exists in the clients map. If it does, the connection is closed, removed from the clients
// start 是 ClientManager 结构体的方法。它启动一个 for 循环，该循环连续监听多个通道上的传入事件。
// 当在注册通道上收到连接时，将该连接添加到 clients 映射中，并向客户端发送指示新套接字连接的 JSON 消息。
// 当在注销通道上收到连接时，将检查该连接是否存在于 clients 映射中。如果是，则关闭该连接，从 clients 中删除该连接
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register: // 当在注册通道上收到连接时，赋值给 conn
			manager.clients[conn] = true                                                  // 将该连接添加到 clients 映射中
			jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has connected."}) // 创建一个 JSON 消息，指示新的套接字连接
			manager.send(jsonMessage, conn)                                               // 向所有客户端发送指示新套接字连接的 JSON 消息
		case conn := <-manager.unregister: // 当在注销通道上收到连接时，赋值给 conn
			if _, ok := manager.clients[conn]; ok { // 检查该连接是否存在于 clients 映射中
				close(conn.send)                                                                  // 如果存在，则关闭该连接
				delete(manager.clients, conn)                                                     // 从 clients 中删除该连接
				jsonMessage, _ := json.Marshal(&Message{Content: "/A s ocket has disconnected."}) // 创建一个 JSON 消息，指示套接字连接已断开
				manager.send(jsonMessage, conn)                                                   // 向所有客户端发送指示套接字连接已断开的 JSON 消息
			}
		case message := <-manager.broadcast: // 当在广播通道上收到消息时，赋值给 message
			for conn := range manager.clients { // 遍历 clients 映射中的所有连接
				select { // 将消息发送到每个连接的 send 通道
				case conn.send <- message: // 如果成功发送消息，则继续循环
				default: // 如果发送失败，则关闭该连接并从 clients 映射中删除该连接
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}
