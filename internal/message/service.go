package message

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

type MessageService struct {
	msgRepo MessageRepository
}

func NewMessageService(msgRepo MessageRepository) *MessageService {
	return &MessageService{
		msgRepo: msgRepo,
	}
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *MessageService) Create(ctx context.Context, sender, recipient, content string) (*Message, error) {
	msg := &Message{
		ID:        generateID(),
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
		CreatedAt: time.Now(),
		Delivered: false,
		Read:      false,
	}

	if err := s.msgRepo.Create(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return msg, nil
}

func (s *MessageService) GetByID(ctx context.Context, id string) (*Message, error) {
	msg, err := s.msgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return msg, nil
}

func (s *MessageService) GetByRecipient(ctx context.Context, recipient string, limit int) ([]*Message, error) {
	msgs, err := s.msgRepo.GetByRecipient(ctx, recipient, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	return msgs, nil
}

func (s *MessageService) GetUndeliveredByRecipient(ctx context.Context, recipient string) ([]*Message, error) {
	msgs, err := s.msgRepo.GetUndeliveredByRecipient(ctx, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to get undelivered messages: %w", err)
	}
	return msgs, nil
}

func (s *MessageService) GetUnreadByRecipient(ctx context.Context, recipient string) ([]*Message, error) {
	msgs, err := s.msgRepo.GetUnreadByRecipient(ctx, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread messages: %w", err)
	}
	return msgs, nil
}

func (s *MessageService) GetMessageHistory(ctx context.Context, user1, user2 string, limit int) ([]*Message, error) {
	msgs, err := s.msgRepo.GetMessageHistory(ctx, user1, user2, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get message history: %w", err)
	}
	return msgs, nil
}

func (s *MessageService) GetUndeliveredMessages(ctx context.Context, recipient string) ([]*Message, error) {
	msgs, err := s.msgRepo.GetUndeliveredMessages(ctx, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to get undelivered messages: %w", err)
	}
	return msgs, nil
}

func (s *MessageService) MarkAsDelivered(ctx context.Context, id string) error {
	if err := s.msgRepo.UpdateDelivered(ctx, id, true); err != nil {
		return fmt.Errorf("failed to mark message as delivered: %w", err)
	}
	return nil
}

func (s *MessageService) MarkAsRead(ctx context.Context, id string) error {
	if err := s.msgRepo.UpdateRead(ctx, id, true); err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}
	return nil
}

func (s *MessageService) Delete(ctx context.Context, id string) error {
	if err := s.msgRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}
