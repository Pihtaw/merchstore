package repository

import (
	"github.com/jmoiron/sqlx"
	"merch-store/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByUserID(userID int) ([]models.Transaction, error)
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	query := `INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, transaction.FromUserID, transaction.ToUserID, transaction.Amount).
		Scan(&transaction.ID, &transaction.CreatedAt)
}

func (r *transactionRepository) GetByUserID(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	// Возвращаем транзакции, где пользователь выступает как отправитель или получатель
	query := `SELECT * FROM transactions WHERE from_user_id=$1 OR to_user_id=$1`
	err := r.db.Select(&transactions, query, userID)
	return transactions, err
}
