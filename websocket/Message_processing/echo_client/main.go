// 命令行 WebSocket 客户端：Dial、WriteMessage、ReadMessage。
package main

import (
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
	text := flag.String("send", `{"op":"echo","text":"hello"}`, "text frame to send (UTF-8)")
	flag.Parse()

	u, err := url.Parse(*addr)
	if err != nil {
		log.Fatal(err)
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read: %v", err)
				return
			}
			log.Printf("recv type=%d: %s", mt, string(message))
		}
	}()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(*text)); err != nil {
		log.Fatalf("write: %v", err)
	}

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("write close: %v", err)
			}
			return
		case <-time.After(30 * time.Second):
			log.Println("timeout")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("write close: %v", err)
			}
			return
		}
	}
}
