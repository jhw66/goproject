// 与 encryption/server 联调：通过 WSS 连接；若服务端启用 WS_ENC_KEY，则帧为 AES-GCM 封装。
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const wssURL = "wss://127.0.0.1:9102/ws"

type payloadCipher struct {
	gcm cipher.AEAD
}

func newPayloadCipherFromEnv() (*payloadCipher, error) {
	b64 := os.Getenv("WS_ENC_KEY")
	if b64 == "" {
		return nil, nil
	}
	key, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("WS_ENC_KEY must decode to 32 bytes")
	}
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(b)
	if err != nil {
		return nil, err
	}
	return &payloadCipher{gcm: g}, nil
}

func (p *payloadCipher) seal(plain []byte) ([]byte, error) {
	if p == nil {
		return plain, nil
	}
	nonce := make([]byte, p.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return p.gcm.Seal(nonce, nonce, plain, nil), nil
}

func (p *payloadCipher) open(data []byte) ([]byte, error) {
	if p == nil {
		return data, nil
	}
	ns := p.gcm.NonceSize()
	if len(data) < ns {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := data[:ns], data[ns:]
	return p.gcm.Open(nil, nonce, ct, nil)
}

func main() {
	pc, err := newPayloadCipherFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	d := websocket.Dialer{
		HandshakeTimeout: 8 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 仅演示自签名证书；生产应校验服务端证书与域名。
		},
	}
	c, res, err := d.Dial(wssURL, nil)
	if err != nil {
		log.Fatalf("dial: %v http=%v", err, res)
	}
	defer c.Close()

	msg := "hello over WSS"
	if pc != nil {
		msg += " + AES-GCM"
	}
	plain := []byte(msg)
	out, err := pc.seal(plain)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.WriteMessage(websocket.TextMessage, out); err != nil {
		log.Fatal(err)
	}
	mt, back, err := c.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	clear, err := pc.open(back)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("recv (%d): %s", mt, clear)
}
