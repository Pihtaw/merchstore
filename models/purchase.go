package models

import "time"

type Purchase struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	MerchID     int       `db:"merch_id" json:"merch_id"`
	PurchasedAt time.Time `db:"purchased_at" json:"purchased_at"`
}
