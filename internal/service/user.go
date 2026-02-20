package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ArchDevs/radix/internal/model"
	"github.com/ArchDevs/radix/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, address string, publicKey []byte) (*model.User, error)
	GetUser(ctx context.Context, address string) (*model.User, error)
	UpdatePublicKey(ctx context.Context, address string, publicKey []byte) error
	DeleteUser(ctx context.Context, address string) error
}

type UserStore struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserStore{userRepo: userRepo}
}

func (s *UserStore) CreateUser(ctx context.Context, address string, publicKey []byte) (*model.User, error) {
	user := &model.User{
		Address:   address,
		PublicKey: publicKey,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserStore) GetUser(ctx context.Context, address string) (*model.User, error) {
	user, err := s.userRepo.GetByAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserStore) UpdatePublicKey(ctx context.Context, address string, publicKey []byte) error {
	if err := s.userRepo.UpdatePublicKey(ctx, address, publicKey); err != nil {
		return fmt.Errorf("failed to update public key: %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(ctx context.Context, address string) error {
	if err := s.userRepo.Delete(ctx, address); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
