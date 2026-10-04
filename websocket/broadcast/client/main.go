// 命令行客户端：连接广播服务，发送一行文本并打印收到的广播。
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			mt, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read: %v", err)
				return
			}
			log.Printf("recv type=%d: %s", mt, string(p))
		}
	}()

	log.Println("输入一行后回车发送，Ctrl+C 退出")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			log.Fatalf("write: %v", err)
		}

		if scanner.Err() != nil {
			log.Printf("stdin: %v", err)
			break
		}
	}

	select {
	case <-interrupt:
	default:
	}

	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
