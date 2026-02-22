package message

import "time"

type Message struct {
	ID        string    `db:"id" json:"id"`
	Sender    string    `db:"sender" json:"sender"`
	Recipient string    `db:"recipient" json:"recipient"`
	Content   string    `db:"content" json:"content"` // encrypted base64
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Delivered bool      `db:"delivered" json:"delivered"`
	Read      bool      `db:"read" json:"read"`
}
