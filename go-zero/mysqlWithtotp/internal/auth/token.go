package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenTypeAccess = "access"
	TokenTypePre2FA = "pre_2fa"
)

type Claims struct {
	UserID    int64  `json:"userId"`
	Username  string `json:"username"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

type Subject struct {
	UserID   int64
	Username string
}

func GenerateToken(secret string, expireSeconds int64, userID int64, username, tokenType string) (string, int64, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return "", 0, errors.New("token secret is required")
	}
	if expireSeconds <= 0 {
		return "", 0, errors.New("token expire must be positive")
	}

	expiresAt := time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()
	claims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return signed, expiresAt, nil
}

func ParseToken(secret, tokenString, expectedType string) (*Subject, error) {
	secret = strings.TrimSpace(secret)
	tokenString = strings.TrimSpace(tokenString)
	if secret == "" || tokenString == "" {
		return nil, errors.New("token secret and token are required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if expectedType != "" && claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return &Subject{
		UserID:   claims.UserID,
		Username: claims.Username,
	}, nil
}
