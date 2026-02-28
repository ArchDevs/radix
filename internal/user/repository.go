package user

import (
	"context"

	"github.com/ArchDevs/radix/internal/database"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByAddress(ctx context.Context, address string) (*User, error)
	UpdatePublicKey(ctx context.Context, address string, publicKey []byte) error
	UpdateUsername(ctx context.Context, address string, username string) error
	Delete(ctx context.Context, address string) error
	Exists(ctx context.Context, address string) (bool, error)
	Search(ctx context.Context, query string, limit int) ([]*User, error)
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
	query := `SELECT address, username, display_name, public_key, created_at FROM users WHERE address = ?`
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

func (r *UserStore) UpdateUsername(ctx context.Context, address string, username string) error {
	query := `UPDATE users SET username = ? WHERE address = ?`
	_, err := r.db.ExecContext(ctx, query, username, address)
	return err
}

func (r *UserStore) Delete(ctx context.Context, address string) error {
	query := `DELETE FROM users WHERE address = ?`
	_, err := r.db.ExecContext(ctx, query, address)
	return err
}

func (r *UserStore) Exists(ctx context.Context, address string) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE address = ?
		)
	`, address).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *UserStore) Search(ctx context.Context, query string, limit int) ([]*User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT address, username, display_name, created_at
		FROM users
		WHERE allow_search = 1
		AND (
			username LIKE ? COLLATE NOCASE
			OR address LIKE ?
		)
		LIMIT ?
	`, query+"%", query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Address, &u.Username, &u.DisplayName, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, rows.Err()
}
