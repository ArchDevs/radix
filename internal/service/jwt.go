package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{secret: []byte(secret), ttl: ttl}
}

func (s *JWTService) Generate(address string) (string, error) {
	claims := jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(s.ttl).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *JWTService) Parse(token string) (string, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	address, ok := claims["address"].(string)
	if !ok {
		return "", fmt.Errorf("address not found")
	}

	return address, nil
}
