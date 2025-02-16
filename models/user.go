package models

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Balance      int       `db:"balance" json:"balance"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
