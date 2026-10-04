package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type outgoingMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type joinRequest struct {
	Password string `json:"password"`
}

func main() {
	addr := flag.String("addr", "http://localhost:3333", "chat server address, e.g. http://localhost:3333")
	room := flag.String("room", "public-room", "room id")
	token := flag.String("token", "", "jwt token")
	password := flag.String("password", "", "optional private room password")
	flag.Parse()

	if strings.TrimSpace(*token) == "" {
		log.Fatal("token is required")
	}

	httpBase := strings.TrimRight(*addr, "/")
	if strings.TrimSpace(*password) != "" {
		if err := joinRoom(httpBase, *room, *token, *password); err != nil {
			log.Fatalf("join private room failed: %v", err)
		}
	}

	wsURL, err := buildWSURL(httpBase, *room, *token)
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatalf("websocket connect failed: %v", err)
	}
	defer conn.Close()

	fmt.Println("connected. type message and press Enter, /quit to exit.")

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read failed: %v", err)
				return
			}
			fmt.Println(string(msg))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "/quit" {
			fmt.Println("bye")
			return
		}

		payload, err := json.Marshal(outgoingMessage{
			Type:    "message",
			Content: line,
		})
		if err != nil {
			log.Printf("marshal failed: %v", err)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("send failed: %v", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("stdin read failed: %v", err)
	}
}

func joinRoom(addr, room, token, password string) error {
	endpoint := fmt.Sprintf("%s/rooms/%s/join?token=%s",
		addr,
		url.PathEscape(room),
		url.QueryEscape(token),
	)
	body, err := json.Marshal(joinRequest{Password: password})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	return nil
}

func buildWSURL(httpAddr, room, token string) (string, error) {
	parsed, err := url.Parse(httpAddr)
	if err != nil {
		return "", err
	}
	scheme := "ws"
	if parsed.Scheme == "https" {
		scheme = "wss"
	}
	q := url.Values{}
	q.Set("roomId", room)
	q.Set("token", token)
	return fmt.Sprintf("%s://%s/ws?%s", scheme, parsed.Host, q.Encode()), nil
}
