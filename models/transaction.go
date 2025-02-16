package models

import "time"

type Transaction struct {
	ID         int       `db:"id" json:"id"`
	FromUserID int       `db:"from_user_id" json:"from_user_id"`
	ToUserID   int       `db:"to_user_id" json:"to_user_id"`
	Amount     int       `db:"amount" json:"amount"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
