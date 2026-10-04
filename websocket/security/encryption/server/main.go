// 示例二：WebSocket 数据加密。
//
// 1) 传输层：使用 WSS（TLS）加密链路，防止窃听与篡改。
// 2) 应用层（可选）：在已 TLS 保护的前提下，仍可对消息再做 AES-GCM 封装（端到端、合规分区等场景）。
//
// 环境变量 WS_ENC_KEY：32 字节的 base64（标准编码），服务端与客户端须一致。未设置则仅演示明文帧（仍走 WSS）。
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	listenAddr = ":9102"
	wsPath     = "/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// devTLSCert 生成自签名证书，仅用于本地演示；生产环境使用正规 CA 或内网 PKI 签发的证书。
func devTLSCert() (tls.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}
	tpl := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{Organization: []string{"websocket-security-demo"}},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:       []string{"localhost", "127.0.0.1"},
	}
	der, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return tls.X509KeyPair(certPEM, keyPEM)
}

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
		return nil, errors.New("WS_ENC_KEY must decode to 32 bytes for AES-256")
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

func serveWS(pc *payloadCipher) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}
		defer conn.Close()

		for {
			mt, payload, err := conn.ReadMessage()
			if err != nil {
				break
			}
			plain, err := pc.open(payload)
			if err != nil {
				log.Printf("decrypt: %v", err)
				break
			}
			out, err := pc.seal(plain)
			if err != nil {
				log.Printf("encrypt: %v", err)
				break
			}
			if err := conn.WriteMessage(mt, out); err != nil {
				break
			}
		}
	}
}

func main() {
	cert, err := devTLSCert()
	if err != nil {
		log.Fatal(err)
	}
	pc, err := newPayloadCipherFromEnv()
	if err != nil {
		log.Fatalf("payload cipher: %v", err)
	}
	if pc == nil {
		log.Print("WS_ENC_KEY unset: echoing plaintext frames over WSS only")
	} else {
		log.Print("WS_ENC_KEY set: frames are AES-GCM sealed (nonce||ciphertext)")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET(wsPath, serveWS(pc))

	srv := &http.Server{
		Addr:              listenAddr,
		Handler:           r,
		ReadHeaderTimeout: 8 * time.Second,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		},
	}
	log.Printf("encryption demo: listen wss://127.0.0.1%s%s (trust dev cert on client)", listenAddr, wsPath)
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
