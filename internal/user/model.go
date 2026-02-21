package user

import (
	"database/sql"
	"time"
)

type User struct {
	Address     string         `db:"address"` // ID
	PublicKey   []byte         `db:"public_key"`
	Username    sql.NullString `db:"username"`
	DisplayName sql.NullString `db:"display_name"`
	CreatedAt   time.Time      `db:"created_at"`
}
