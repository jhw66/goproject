package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{
		"error": msg,
	})
}

func ReadJSON(r *http.Request, out interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func AuthFromRequest(r *http.Request, cfg AppConfig) (*UserClaims, error) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			token = strings.TrimSpace(auth[7:])
		}
	}
	return ParseUserClaims(token, cfg.JWTSecret)
}

func ParseLimit(s string, fallback int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return fallback
	}
	if n > 200 {
		return 200
	}
	return n
}

func ParseRoomID(path string) string {
	const prefix = "/rooms/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	rest := strings.TrimPrefix(path, prefix)
	parts := strings.Split(rest, "/")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}

func ParseRoomSuffix(path string) string {
	const prefix = "/rooms/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	rest := strings.TrimPrefix(path, prefix)
	parts := strings.Split(rest, "/")
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func NewRoomID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
