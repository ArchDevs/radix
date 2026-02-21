package auth

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/ArchDevs/radix/internal/user"
)

type AuthService struct {
	userSvc user.UserService
}

func NewAuthService(userSvc user.UserService) *AuthService {
	return &AuthService{userSvc: userSvc}
}

func (s *AuthService) Register(ctx context.Context, address string, pubKeyBase64 string) (*user.User, error) {
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
