package user

import (
	"context"

	"github.com/ArchDevs/radix/internal/database"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByAddress(ctx context.Context, address string) (*User, error)
	UpdatePublicKey(ctx context.Context, address string, publicKey []byte) error
	Delete(ctx context.Context, address string) error
}

type UserStore struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &UserStore{db: db}
}

func (r *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (address, public_key, created_at) VALUES (:address, :public_key, :created_at)`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserStore) GetByAddress(ctx context.Context, address string) (*User, error) {
	query := `SELECT address, public_key, created_at FROM users WHERE address = ?`
	var user User
	err := r.db.GetContext(ctx, &user, query, address)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserStore) UpdatePublicKey(ctx context.Context, address string, publicKey []byte) error {
	query := `UPDATE users SET public_key = ? WHERE address = ?`
	_, err := r.db.ExecContext(ctx, query, publicKey, address)
	return err
}

func (r *UserStore) Delete(ctx context.Context, address string) error {
	query := `DELETE FROM users WHERE address = ?`
	_, err := r.db.ExecContext(ctx, query, address)
	return err
}
