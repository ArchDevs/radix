package model

import "time"

type Message struct {
	ID         string    `db:"id"`
	Sender     string    `db:"sender"`
	Recipient  string    `db:"recipient"`
	Ciphertext []byte    `db:"ciphertext"`
	CreatedAt  time.Time `db:"created_at"`
	Delivered  bool      `db:"delivered"`
}
