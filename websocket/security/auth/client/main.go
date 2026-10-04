// 与 auth/server 联调：先 POST /login 取 JWT，再带 token 连接 WebSocket。
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const serverBase = "http://127.0.0.1:8888"

type loginResp struct {
	Token string `json:"token"`
}

func main() {
	body := bytes.NewBufferString(`{"user":"demo"}`)
	req, err := http.NewRequest(http.MethodPost, serverBase+"/login", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("login: %s", resp.Status)
	}
	var lr loginResp
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		log.Fatal(err)
	}
	if lr.Token == "" {
		log.Fatal("no token in response")
	}

	u, err := url.Parse(serverBase)
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "ws"
	u.Path = "/ws"
	q := u.Query()
	q.Set("token", lr.Token)
	u.RawQuery = q.Encode()

	d := websocket.Dialer{HandshakeTimeout: 8 * time.Second}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+lr.Token)
	header.Set("Sec-WebSocket-Protocol", "bearer."+lr.Token)

	c, res, err := d.Dial(u.String(), header)
	if err != nil {
		log.Fatalf("dial: %v (http=%v)", err, res)
	}
	defer c.Close()

	mt, msg, err := c.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("server: %s (%d)", msg, mt)

	if err := c.WriteMessage(websocket.TextMessage, []byte("ping from client")); err != nil {
		log.Fatal(err)
	}
	mt, msg, err = c.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("echo: %s", msg)
}
