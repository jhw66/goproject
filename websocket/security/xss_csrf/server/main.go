// 示例三：与 WebSocket 相关的 XSS / CSRF 防护思路。
//
// CSRF（跨站请求伪造）：
//   - 恶意站点诱导浏览器向你的站点发起已登录请求。对 WebSocket 而言，若握手会携带会话 Cookie，
//     且未校验 Origin，攻击者页面可能建立到你域名的 WebSocket。
//   - 缓解：严格校验 Origin；Cookie 使用 SameSite（Lax/Strict）；敏感操作用二次确认或 CSRF Token。
//
// 本示例：在握手时校验 Origin 白名单，并要求查询参数 csrf 与名为 ws_csrf 的 Cookie 一致（双重提交）。
//
// XSS（跨站脚本）：
//   - 若将服务端消息直接写入 innerHTML，恶意脚本可能执行。应用 textContent / 模板转义 / CSP。
//   - 服务端对输出做编码（如 JSON）、对 HTML 页面设置 Content-Security-Policy。
package main

import (
	"crypto/rand"
	"encoding/hex"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	listenAddr = ":9103"
	wsPath     = "/ws"
	csrfCookie = "ws_csrf"
)

// allowedOrigins 为 WebSocket Origin 白名单；可用环境变量 WS_ALLOWED_ORIGINS（逗号分隔）覆盖。
func loadAllowedOrigins() map[string]struct{} {
	raw := os.Getenv("WS_ALLOWED_ORIGINS")
	if raw == "" {
		raw = "http://127.0.0.1:9103,http://localhost:9103"
	}
	m := make(map[string]struct{})
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			m[o] = struct{}{}
		}
	}
	return m
}

func newUpgrader(allowed map[string]struct{}) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return false
			}
			_, ok := allowed[origin]
			return ok
		},
	}
}

func issueCSRFCookie(c *gin.Context) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	token := hex.EncodeToString(b)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     csrfCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false, // 演示页需读取以拼 WebSocket URL；生产可改为仅服务端比对 Cookie 与 POST body。
		SameSite: http.SameSiteStrictMode,
	})
	c.JSON(http.StatusOK, gin.H{"csrf": token})
}

func servePage(c *gin.Context) {
	c.Header("Content-Security-Policy", "default-src 'self'; base-uri 'none'; frame-ancestors 'none'")
	// 演示：内联脚本需 unsafe-inline；生产建议全部走静态文件 + nonce/sha256。
	c.Header("Content-Type", "text/html; charset=utf-8")
	const page = `<!doctype html>
<meta charset="utf-8">
<title>WS XSS/CSRF demo</title>
<pre id="log"></pre>
<script>
(async () => {
  const log = (m) => { const el = document.getElementById('log'); el.textContent += m + "\n"; };
  const r = await fetch('/api/bootstrap', { credentials: 'same-origin' });
  const j = await r.json();
  const csrf = j.csrf;
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
  const ws = new WebSocket(proto + '//' + location.host + '/ws?csrf=' + encodeURIComponent(csrf));
  ws.onopen = () => log('open');
  ws.onmessage = (ev) => {
    // 使用 textContent 展示，避免把消息当 HTML 解析（缓解 XSS）。
    const safe = document.createElement('div');
    safe.textContent = ev.data;
    log('message: ' + safe.textContent);
  };
  ws.onerror = () => log('error');
})();
</script>
`
	c.String(http.StatusOK, page)
}

// safeEchoJSON 将用户输入作为 JSON 字符串写出，避免拼接到 HTML（encoding/json 会转义）。
func safeEchoJSON(c *gin.Context) {
	var body struct{ Msg string `json:"msg"` }
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"echo": body.Msg})
}

// unsafeHTMLDemo 仅作对照：若把用户输入嵌入 HTML 必须转义，否则存在 XSS。
func unsafeHTMLDemo(c *gin.Context) {
	q := c.Query("q")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, "<p>escaped: "+html.EscapeString(q)+"</p>")
}

func serveWS(up *websocket.Upgrader) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(csrfCookie)
		if err != nil || cookie.Value == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "missing csrf cookie"})
			return
		}
		if c.Query("csrf") != cookie.Value {
			c.JSON(http.StatusForbidden, gin.H{"error": "csrf mismatch"})
			return
		}

		conn, err := up.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}
		defer conn.Close()

		welcome := "<b>hello</b> — 若前端用 innerHTML 渲染会变成粗体；本示例客户端用 textContent 显示纯文本。"
		if err := conn.WriteMessage(websocket.TextMessage, []byte(welcome)); err != nil {
			return
		}
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
}

func main() {
	allowed := loadAllowedOrigins()
	up := newUpgrader(allowed)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/", servePage)
	r.GET("/api/bootstrap", issueCSRFCookie)
	r.POST("/api/echo", safeEchoJSON)
	r.GET("/unsafe-html-demo", unsafeHTMLDemo)
	r.GET(wsPath, serveWS(&up))

	log.Printf("xss/csrf demo: open http://127.0.0.1%s  (allowed Origins: %v)", listenAddr, allowed)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}
