package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// UserClaims carries JWT identity information.
type UserClaims struct {
	UserID   string
	Nickname string
}

type tokenClaims struct {
	UserID   string `json:"userId"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

func ParseUserClaims(tokenString, secret string) (*UserClaims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	parsedClaim, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedClaim.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedClaim.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	if strings.TrimSpace(claims.UserID) == "" {
		return nil, errors.New("missing userId in token")
	}
	name := strings.TrimSpace(claims.Nickname)
	if name == "" {
		name = claims.UserID
	}

	return &UserClaims{
		UserID:   claims.UserID,
		Nickname: name,
	}, nil
}
