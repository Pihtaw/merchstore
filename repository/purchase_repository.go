package repository

import (
	"github.com/jmoiron/sqlx"
	"merch-store/models"
)

type PurchaseRepository interface {
	Create(purchase *models.Purchase) error
	GetByUserID(userID int) ([]models.Purchase, error)
}

type purchaseRepository struct {
	db *sqlx.DB
}

func NewPurchaseRepository(db *sqlx.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) Create(purchase *models.Purchase) error {
	query := `INSERT INTO purchases (user_id, merch_id) VALUES ($1, $2) RETURNING id, purchased_at`
	return r.db.QueryRow(query, purchase.UserID, purchase.MerchID).Scan(&purchase.ID, &purchase.PurchasedAt)
}

func (r *purchaseRepository) GetByUserID(userID int) ([]models.Purchase, error) {
	var purchases []models.Purchase
	query := `SELECT * FROM purchases WHERE user_id=$1`
	err := r.db.Select(&purchases, query, userID)
	return purchases, err
}
