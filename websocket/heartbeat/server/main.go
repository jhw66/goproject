// WebSocket 心跳检测示例服务端：周期性向客户端发送 Ping 控制帧，
// 依赖对端返回 Pong；读超时未收到任何数据或 Pong 则关闭连接。
// 所有写操作集中在 writePump，避免与 ReadMessage 并发写同一连接。
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	listenAddr = ":8888"
	wsPath     = "/ws"

	// pongWait：允许「下一次读」的最长等待时间；收到任意消息或 Pong 会刷新。
	pongWait = 9 * time.Second
	// pingPeriod：发送 Ping 的间隔，须小于 pongWait（一般为 pongWait 的 9/10）。
	pingPeriod = (pongWait * 9) / 10

	writeWait = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade: %v", err)
		return
	}

	// 客户端发来 Ping 时由库默认回复 Pong；也可自定义 SetPingHandler。
	send := make(chan []byte, 16)
	go writePump(conn, send)
	readPump(conn, send)
}

func readPump(conn *websocket.Conn, send chan<- []byte) {
	defer close(send)

	conn.SetReadLimit(1024)
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("set read deadline: %v", err)
		return
	}
	conn.SetPongHandler(func(appData string) error {
		log.Println("收到 Pong，刷新读超时,:", appData)
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read: %v", err)
			}
			return
		}
		if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return
		}
		log.Printf("收到消息 type=%d: %s", mt, string(p))
		if mt == websocket.TextMessage {
			select {
			case send <- append([]byte(nil), p...):
			default:
				log.Println("send 队列满，丢弃回显")
			}
		}
	}
}

func writePump(conn *websocket.Conn, send <-chan []byte) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-send:
			if !ok {
				deadline := time.Now().Add(writeWait)
				_ = conn.SetWriteDeadline(deadline)
				_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("write text: %v", err)
				return
			}
		case <-ticker.C:
			//_ = conn.SetWriteDeadline(time.Now().Add(writeWait)) 不需要，因为WriteControl会自动设置写超时
			if err := conn.WriteControl(websocket.PingMessage, []byte("srv-ping"), time.Now().Add(writeWait)); err != nil {
				log.Printf("write ping: %v", err)
				return
			}
			log.Println("已发送 Ping")
		}
	}
}

func main() {
	r := gin.Default()
	r.GET(wsPath, serveWS)
	log.Printf("heartbeat server listening on %s  GET %s", listenAddr, wsPath)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
