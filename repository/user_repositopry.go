package repository

import (
	"github.com/jmoiron/sqlx"
	"merch-store/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	UpdateBalance(userID int, newBalance int) error
	GetByID(id int) (*models.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, password_hash, balance) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, user.Username, user.PasswordHash, user.Balance).Scan(&user.ID)
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE username=$1`
	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id=$1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateBalance(userID int, newBalance int) error {
	query := `UPDATE users SET balance=$1 WHERE id=$2`
	_, err := r.db.Exec(query, newBalance, userID)
	return err
}
