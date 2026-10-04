package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader：用于将HTTP连接升级为WebSocket连接。CheckOrigin函数用于检查请求的来源，这里我们简单地返回true，允许所有来源的连接。
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handleWebSocket：处理WebSocket请求的函数。它首先将HTTP连接升级为WebSocket连接，然后进入一个循环，不断读取客户端发送的消息并将其回显给客户端。
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}

// main：启动Gin服务器，并将/ws路径映射到handleWebSocket函数。
func main() {
	r := gin.Default()
	r.GET("/ws", handleWebSocket)
	r.Run(":8888")
}
