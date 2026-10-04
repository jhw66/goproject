// 示例一：WebSocket 握手前身份验证（JWT）。
//
// 要点：
//   - 在 Upgrade 之前校验令牌；失败则返回 401，不建立 WebSocket。
//   - 令牌可来自查询参数、Authorization 头或 Sec-WebSocket-Protocol（部分代理对自定义头不友好时可选用其一）。
//   - 生产环境应使用 HTTPS/WSS、短过期时间、密钥轮换与 claims 校验（aud、iss 等）。
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

const (
	listenAddr = ":8888"
	wsPath     = "/ws"
	jwtSecret  = "dev-only-change-me-use-32bytes-min____"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func parseToken(r *http.Request) (string, error) {
	if t := r.URL.Query().Get("token"); t != "" {
		fmt.Println("token from query:", t)
		return t, nil
	}
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(h), "bearer ") {
		fmt.Println("token from header:", strings.TrimSpace(h[7:]))
		return strings.TrimSpace(h[7:]), nil
	}
	// 部分客户端通过子协议传递：Sec-WebSocket-Protocol: bearer, <jwt>
	if p := r.Header.Get("Sec-WebSocket-Protocol"); p != "" {
		parts := strings.Split(p, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		for _, part := range parts {
			if strings.HasPrefix(part, "bearer.") {
				fmt.Println("token from Sec-WebSocket-Protocol:", strings.TrimPrefix(part, "bearer."))
				return strings.TrimPrefix(part, "bearer."), nil
			}
		}
	}
	return "", errors.New("missing token")
}

func validateJWT(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func issueDemoToken(c *gin.Context) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   "demo-user",
		Issuer:    "security-demo",
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sign failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": s, "ws_url_query_example": wsPath + "?token=" + s})
}

func serveWS(c *gin.Context) {
	raw, err := parseToken(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if _, err := validateJWT(raw); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade: %v", err)
		return
	}
	defer conn.Close()

	_ = conn.WriteMessage(websocket.TextMessage, []byte("authenticated websocket ok"))
	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := conn.WriteMessage(mt, p); err != nil {
			break
		}
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/login", issueDemoToken)
	r.GET(wsPath, serveWS)

	log.Printf("auth demo: listen %s  POST /login  GET %s?token=<jwt>", listenAddr, wsPath)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
