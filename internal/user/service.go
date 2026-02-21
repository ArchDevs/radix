package user

import (
	"context"
	"fmt"
	"time"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, address string, publicKey []byte) (*User, error) {
	user := &User{
		Address:   address,
		PublicKey: publicKey,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, address string) (*User, error) {
	user, err := s.userRepo.GetByAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdatePublicKey(ctx context.Context, address string, publicKey []byte) error {
	if err := s.userRepo.UpdatePublicKey(ctx, address, publicKey); err != nil {
		return fmt.Errorf("failed to update public key: %w", err)
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, address string) error {
	if err := s.userRepo.Delete(ctx, address); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
