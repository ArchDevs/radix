package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/ArchDevs/radix/internal/model"
)

type AuthService struct {
	userSvc UserService
}

func NewAuthService(userSvc UserService) *AuthService {
	return &AuthService{userSvc: userSvc}
}

func (s *AuthService) Register(ctx context.Context, address string, pubKeyBase64 string) (*model.User, error) {
	pubKey, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid public key encoding: %w", err)
	}

	user, err := s.userSvc.CreateUser(ctx, address, pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}
