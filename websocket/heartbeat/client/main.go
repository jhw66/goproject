// 命令行客户端：连接心跳示例服务，打印收到的消息与 Ping/Pong 相关日志。
// 本客户端同样用单独写协程发 Ping，读协程处理 Pong 与读超时。
package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 9 * time.Second
	pingPeriod = (pongWait * 9) / 10
	writeWait  = 10 * time.Second
)

func main() {
	addr := flag.String("url", "ws://127.0.0.1:8888/ws", "WebSocket URL")
	flag.Parse()

	u, err := url.Parse(*addr)
	if err != nil {
		log.Fatal(err)
	}

	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	send := make(chan []byte, 16)
	done := make(chan struct{})

	go writePump(conn, send, done)
	go readPump(conn, send, done)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		log.Println("输入一行后回车发送文本；Ctrl+C 退出")
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			line := sc.Text()
			select {
			case send <- []byte(line):
			case <-done:
				return
			}
		}
		if err := sc.Err(); err != nil {
			log.Printf("stdin: %v", err)
		}
	}()

	select {
	case <-interrupt:
	case <-done:
	}
	// 仅关闭连接，避免与 writePump 并发 Write。
	_ = conn.Close()
}

func readPump(conn *websocket.Conn, send chan<- []byte, done chan<- struct{}) {
	defer func() { close(send); close(done) }()

	conn.SetReadLimit(1024)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		log.Println("收到服务端 Pong 刷新读超时,:", appData)
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read: %v", err)
			return
		}
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Printf("recv type=%d: %s", mt, string(p))
	}
}

func writePump(conn *websocket.Conn, send <-chan []byte, done <-chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case msg, ok := <-send:
			if !ok {
				return
			}
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("write: %v", err)
				_ = conn.Close()
				return
			}
		case <-ticker.C:
			//_ = conn.SetWriteDeadline(time.Now().Add(writeWait)) 不需要，因为WriteControl会自动设置写超时
			if err := conn.WriteControl(websocket.PingMessage, []byte("cli-ping"), time.Now().Add(writeWait)); err != nil {
				log.Printf("ping: %v", err)
				_ = conn.Close()
				return
			}
			log.Println("已发送 Ping（至服务端）")
		}
	}
}
