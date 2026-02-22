package message

import (
	"context"

	"github.com/ArchDevs/radix/internal/database"
)

type MessageRepository interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, id string) (*Message, error)
	GetByRecipient(ctx context.Context, recipient string, limit int) ([]*Message, error)
	GetUndeliveredByRecipient(ctx context.Context, recipient string) ([]*Message, error)
	GetUnreadByRecipient(ctx context.Context, recipient string) ([]*Message, error)
	GetMessageHistory(ctx context.Context, user1, user2 string, limit int) ([]*Message, error)
	GetUndeliveredMessages(ctx context.Context, recipient string) ([]*Message, error)
	UpdateDelivered(ctx context.Context, id string, delivered bool) error
	UpdateRead(ctx context.Context, id string, read bool) error
	Delete(ctx context.Context, id string) error
}

type MessageStore struct {
	db *database.DB
}

func NewMessageRepository(db *database.DB) MessageRepository {
	return &MessageStore{db: db}
}

func (r *MessageStore) Create(ctx context.Context, message *Message) error {
	query := `INSERT INTO messages (id, sender, recipient, content, created_at, delivered, read) 
			  VALUES (:id, :sender, :recipient, :content, :created_at, :delivered, :read)`
	_, err := r.db.NamedExecContext(ctx, query, message)
	return err
}

func (r *MessageStore) GetByID(ctx context.Context, id string) (*Message, error) {
	query := `SELECT id, sender, recipient, content, created_at, delivered, read FROM messages WHERE id = ?`
	var message Message
	err := r.db.GetContext(ctx, &message, query, id)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageStore) GetByRecipient(ctx context.Context, recipient string, limit int) ([]*Message, error) {
	query := `SELECT id, sender, recipient, content, created_at, delivered, read 
			  FROM messages WHERE recipient = ? ORDER BY created_at DESC LIMIT ?`
	var messages []*Message
	err := r.db.SelectContext(ctx, &messages, query, recipient, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageStore) GetUndeliveredByRecipient(ctx context.Context, recipient string) ([]*Message, error) {
	query := `SELECT id, sender, recipient, content, created_at, delivered, read 
			  FROM messages WHERE recipient = ? AND delivered = false ORDER BY created_at ASC`
	var messages []*Message
	err := r.db.SelectContext(ctx, &messages, query, recipient)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageStore) GetUnreadByRecipient(ctx context.Context, recipient string) ([]*Message, error) {
	query := `SELECT id, sender, recipient, content, created_at, delivered, read 
			  FROM messages WHERE recipient = ? AND read = false ORDER BY created_at ASC`
	var messages []*Message
	err := r.db.SelectContext(ctx, &messages, query, recipient)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageStore) GetMessageHistory(ctx context.Context, user1, user2 string, limit int) ([]*Message, error) {
	query := `SELECT id, sender, recipient, content, created_at, delivered, read 
			  FROM messages 
			  WHERE (sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)
			  ORDER BY created_at DESC
			  LIMIT ?`
	var messages []*Message
	err := r.db.SelectContext(ctx, &messages, query, user1, user2, user2, user1, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageStore) GetUndeliveredMessages(ctx context.Context, recipient string) ([]*Message, error) {
	query := `SELECT id, sender, recipient, content, created_at, delivered, read 
			  FROM messages WHERE recipient = ? AND delivered = false ORDER BY created_at ASC`
	var messages []*Message
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

func (r *MessageStore) UpdateRead(ctx context.Context, id string, read bool) error {
	query := `UPDATE messages SET read = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, read, id)
	return err
}

func (r *MessageStore) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
