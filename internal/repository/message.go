package repository

import (
	"context"

	"github.com/ArchDevs/radix/internal/database"
	"github.com/ArchDevs/radix/internal/model"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	GetByID(ctx context.Context, id string) (*model.Message, error)
	GetByRecipient(ctx context.Context, recipient string, limit int) ([]model.Message, error)
	GetUndeliveredByRecipient(ctx context.Context, recipient string) ([]model.Message, error)
	UpdateDelivered(ctx context.Context, id string, delivered bool) error
	Delete(ctx context.Context, id string) error
}

type MessageStore struct {
	db *database.DB
}

func NewMessageRepository(db *database.DB) MessageRepository {
	return &MessageStore{db: db}
}

func (r *MessageStore) Create(ctx context.Context, message *model.Message) error {
	query := `INSERT INTO messages (id, sender, recipient, ciphertext, created_at, delivered) 
			  VALUES (:id, :sender, :recipient, :ciphertext, :created_at, :delivered)`
	_, err := r.db.NamedExecContext(ctx, query, message)
	return err
}

func (r *MessageStore) GetByID(ctx context.Context, id string) (*model.Message, error) {
	query := `SELECT id, sender, recipient, ciphertext, created_at, delivered FROM messages WHERE id = ?`
	var message model.Message
	err := r.db.GetContext(ctx, &message, query, id)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageStore) GetByRecipient(ctx context.Context, recipient string, limit int) ([]model.Message, error) {
	query := `SELECT id, sender, recipient, ciphertext, created_at, delivered 
			  FROM messages WHERE recipient = ? ORDER BY created_at DESC LIMIT ?`
	var messages []model.Message
	err := r.db.SelectContext(ctx, &messages, query, recipient, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageStore) GetUndeliveredByRecipient(ctx context.Context, recipient string) ([]model.Message, error) {
	query := `SELECT id, sender, recipient, ciphertext, created_at, delivered 
			  FROM messages WHERE recipient = ? AND delivered = false ORDER BY created_at ASC`
	var messages []model.Message
	err := r.db.SelectContext(ctx, &messages, query, recipient)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageStore) UpdateDelivered(ctx context.Context, id string, delivered bool) error {
	query := `UPDATE messages SET delivered = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, delivered, id)
	return err
}

func (r *MessageStore) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
