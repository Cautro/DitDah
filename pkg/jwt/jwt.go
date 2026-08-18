package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"crypto/rand"
	"encoding/base64"
)

type TokenManager struct {
	secretKey []byte
	ttl       time.Duration
}

type Claims struct {
	UserID   int `json:"Id"`
	Username string `json:"Username"`
	jwt.RegisteredClaims
}

func NewTokenManager(secretKey string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(secretKey),
		ttl:       15 * time.Minute,
	}
}

func (m *TokenManager) GenerateAccessToken(userId int, username string) (string, error) {
	claims := &Claims{
		UserID: userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *TokenManager) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	return claims, nil
}

func (m *TokenManager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}