package service

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ArchDevs/radix/internal/model"
	"github.com/ArchDevs/radix/internal/repository"
)

type ChallengeService struct {
	challengeRepo repository.ChallengeRepository

	userSvc UserService
}

func NewChallengeService(challengeRepo repository.ChallengeRepository, userSvc UserService) *ChallengeService {
	return &ChallengeService{
		challengeRepo: challengeRepo,
		userSvc:       userSvc,
	}
}

func (s *ChallengeService) CreateChallenge(ctx context.Context, address string) (string, int64, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", 0, fmt.Errorf("failed to generate nonce: %w", err)
	}
	nonceStr := base64.StdEncoding.EncodeToString(nonce)
	createdAt := time.Now()

	challenge := &model.Challenge{
		Address:   address,
		Nonce:     nonceStr,
		CreatedAt: createdAt,
	}

	if err := s.challengeRepo.Upsert(ctx, challenge); err != nil {
		return "", 0, fmt.Errorf("failed to save challenge: %w", err)
	}

	return nonceStr, createdAt.Unix(), nil
}

func (s *ChallengeService) GetChallenge(ctx context.Context, address string) (*model.Challenge, error) {
	challenge, err := s.challengeRepo.GetByAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	return challenge, nil
}

func (s *ChallengeService) Verify(ctx context.Context, address, nonce, signature string) (bool, error) {
	challenge, err := s.challengeRepo.GetByAddress(ctx, address)
	if err != nil {
		return false, fmt.Errorf("failed to get challenge: %w", err)
	}

	if challenge.Nonce != nonce {
		return false, fmt.Errorf("nonce does not match")
	}

	user, err := s.userSvc.GetUser(ctx, address)
	if err != nil {
		return false, fmt.Errorf("user not found, invalid address")
	}

	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature encoding: %w", err)
	}

	message := []byte(nonce)

	if !ed25519.Verify(user.PublicKey, message, sig) {
		return false, fmt.Errorf("invalid signature")
	}

	// Delete challenge after successful verification
	if err := s.challengeRepo.Delete(ctx, address); err != nil {
		// Log error but don't fail verification
		// Challenge will be cleaned up by cleanup job
		_ = err
	}

	return true, nil
}
