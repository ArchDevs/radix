package model

import "time"

type User struct {
	Address   string    `db:"address"` // ID
	PublicKey []byte    `db:"public_key"`
	CreatedAt time.Time `db:"created_at"`
}
