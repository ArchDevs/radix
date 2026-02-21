package challenge

import (
	"context"
	"time"

	"github.com/ArchDevs/radix/internal/database"
)

type ChallengeRepository interface {
	Create(ctx context.Context, challenge *Challenge) error
	GetByAddress(ctx context.Context, address string) (*Challenge, error)
	Update(ctx context.Context, challenge *Challenge) error
	Upsert(ctx context.Context, challenge *Challenge) error
	Delete(ctx context.Context, address string) error
	DeleteOlderThan(ctx context.Context, duration time.Duration) error
}

type ChallengeStore struct {
	db *database.DB
}

func NewChallengeRepository(db *database.DB) ChallengeRepository {
	return &ChallengeStore{db: db}
}

func (r *ChallengeStore) Create(ctx context.Context, challenge *Challenge) error {
	query := `INSERT INTO challenges (address, nonce, created_at) VALUES (:address, :nonce, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, challenge)
	return err
}

func (r *ChallengeStore) GetByAddress(ctx context.Context, address string) (*Challenge, error) {
	query := `SELECT address, nonce, created_at FROM challenges WHERE address = ?`
	var challenge Challenge
	err := r.db.GetContext(ctx, &challenge, query, address)
	if err != nil {
		return nil, err
	}
	return &challenge, nil
}

func (r *ChallengeStore) Update(ctx context.Context, challenge *Challenge) error {
	query := `UPDATE challenges SET nonce = :nonce, created_at = :created_at WHERE address = :address`
	_, err := r.db.NamedExecContext(ctx, query, challenge)
	return err
}

func (r *ChallengeStore) Upsert(ctx context.Context, challenge *Challenge) error {
	query := `INSERT INTO challenges (address, nonce, created_at) 
			  VALUES (:address, :nonce, :created_at)
			  ON CONFLICT(address) DO UPDATE SET 
			    nonce = excluded.nonce, 
			    created_at = excluded.created_at`
	_, err := r.db.NamedExecContext(ctx, query, challenge)
	return err
}

func (r *ChallengeStore) Delete(ctx context.Context, address string) error {
	query := `DELETE FROM challenges WHERE address = ?`
	_, err := r.db.ExecContext(ctx, query, address)
	return err
}

func (r *ChallengeStore) DeleteOlderThan(ctx context.Context, duration time.Duration) error {
	query := `DELETE FROM challenges WHERE created_at < datetime('now', ?)`
	_, err := r.db.ExecContext(ctx, query, "-"+duration.String())
	return err
}
