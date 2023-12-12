package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

// wsPage is a handler function for upgrading a HTTP connection to a WebSocket connection.
// It upgrades the connection using the websocket.Upgrader and registers the client with the ClientManager.
// It then starts separate goroutines to handle reading and writing on the WebSocket connection.
// wsPage 是一个处理程序函数，用于将 HTTP 连接升级为 WebSocket 连接。
// 它使用 websocket.Upgrader 升级连接并使用 ClientManager 注册客户端。
// 然后，它启动单独的 goroutine 来处理 WebSocket 连接上的读取和写入。
func wsPage(res http.ResponseWriter, req *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	// 注意：这里的 CheckOrigin 属性设置为 true，这是因为我们将在本地运行客户端和服务器，因此不需要检查来源。
	// 如果您将此代码部署到生产环境中，则将相应地设置 CheckOrigin 属性。
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}
	// 将客户端注册到 ClientManager 中
	client := &Client{
		id:     uuid.NewV4().String(), // 生成唯一的 ID
		socket: conn,                  // 刚刚生成的 websocket 连接
		send:   make(chan []byte),     // 用于发送消息的通道，它是一个字节切片，用于存储待发送给这个客户端的信息
	}

	manager.register <- client // 将客户端注册到 ClientManager 中

	go client.read() // 启动单独的 goroutine 来处理 WebSocket 连接上的读取和写入
	go client.write()
}

// main is the entry point for the application.
// It initializes a ClientManager and starts it in a separate goroutine.
// It then sets up a WebSocket route and listens on port 12345 for incoming connections.
// If there is an error starting the server, the function returns.
// main 是应用程序的入口点。
// 它初始化一个 ClientManager 并在单独的 goroutine 中启动它。
// 然后它设置一个 WebSocket 路由，并在端口 12345 上监听传入的连接。
// 如果启动服务器时出现错误，则该函数返回。
func main() {
	fmt.Println("Starting application...")
	go manager.start()             // 初始化一个 ClientManager 并在单独的 goroutine 中启动它。
	http.HandleFunc("/ws", wsPage) // 设置一个 WebSocket 路由，并在端口 12345 上监听传入的连接。
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		return
	}
}
