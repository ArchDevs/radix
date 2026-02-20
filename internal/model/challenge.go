package model

import "time"

type Challenge struct {
	Address   string    `db:"address"`
	Nonce     string    `db:"nonce"`
	CreatedAt time.Time `db:"created_at"`
}
