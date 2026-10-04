// 命令行客户端：连接 redis-sync 服务，stdin 一行即发送；打印广播消息。
// 多开几个终端可验证 Redis Pub/Sub + Hub 广播。
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait  = 10 * time.Second
	writeWait = 10 * time.Second
)

func main() {
	addr := flag.String("url", "ws://127.0.0.1:8888/ws", "WebSocket URL")
	user := flag.String("user", "", "可选，对应服务端 ?user=")
	flag.Parse()

	u, err := url.Parse(*addr)
	if err != nil {
		log.Fatal(err)
	}
	if *user != "" {
		q := u.Query()
		q.Set("user", *user)
		u.RawQuery = q.Encode()
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
		log.Println("输入一行回车发送；Ctrl+C 退出")
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
	_ = conn.Close()
}

func readPump(conn *websocket.Conn, send chan<- []byte, done chan<- struct{}) {
	defer func() { close(send); close(done) }()

	conn.SetReadLimit(4096)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(appData string) error {
		conn.WriteControl(websocket.PongMessage, []byte("pong"), time.Now().Add(writeWait))
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read: %v", err)
			return
		}
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		fmt.Println("client recv message from server")
		log.Printf("recv type=%d: %s", mt, string(p))
	}
}

func writePump(conn *websocket.Conn, send <-chan []byte, done <-chan struct{}) {
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
			fmt.Println("client write message to server")
		}
	}
}
