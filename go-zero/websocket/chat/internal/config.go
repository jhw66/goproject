package internal

import (
	"os"
	"strconv"
)

const (
	defaultHost       = "localhost"
	defaultJWTSecret  = "change-me-in-production"
	defaultMySQLDSN   = "wjhaccount:Wujiahui789@tcp(127.0.0.1:3306)/chat?parseTime=true&charset=utf8mb4"
	defaultRedisAddr  = "127.0.0.1:6379"
	defaultOnlineTTL  = 120
	defaultRedisDB    = 0
	defaultRedisPass  = ""
	defaultHTTPTimeMS = 0
	defaultCPULimit   = 500
	defaultPort       = 3333
)

// AppConfig contains runtime configuration.
type AppConfig struct {
	Host             string
	Port             int
	Timeout          int64
	CpuThreshold     int64
	JWTSecret        string
	MySQLDSN         string
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
	OnlineTTLSeconds int
}

func LoadConfig(port int, timeout, cpu int64) AppConfig {
	if port <= 0 {
		port = envInt("CHAT_PORT", defaultPort)
	}
	if timeout < 0 {
		timeout = defaultHTTPTimeMS
	}
	if cpu <= 0 {
		cpu = defaultCPULimit
	}

	return AppConfig{
		Host:             envOrDefault("CHAT_HOST", defaultHost),
		Port:             port,
		Timeout:          timeout,
		CpuThreshold:     cpu,
		JWTSecret:        envOrDefault("CHAT_JWT_SECRET", defaultJWTSecret),
		MySQLDSN:         envOrDefault("CHAT_MYSQL_DSN", defaultMySQLDSN),
		RedisAddr:        envOrDefault("CHAT_REDIS_ADDR", defaultRedisAddr),
		RedisPassword:    envOrDefault("CHAT_REDIS_PASSWORD", defaultRedisPass),
		RedisDB:          envInt("CHAT_REDIS_DB", defaultRedisDB),
		OnlineTTLSeconds: envInt("CHAT_ONLINE_TTL_SECONDS", defaultOnlineTTL),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return val
}
