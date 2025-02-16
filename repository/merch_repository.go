package repository

import (
	"github.com/jmoiron/sqlx"
	"merch-store/models"
)

type MerchRepository interface {
	GetAll() ([]models.Merch, error)
	GetByID(id int) (*models.Merch, error)
	GetByName(name string) (*models.Merch, error)
}

type merchRepository struct {
	db *sqlx.DB
}

func NewMerchRepository(db *sqlx.DB) MerchRepository {
	return &merchRepository{db: db}
}

func (r *merchRepository) GetAll() ([]models.Merch, error) {
	var merchs []models.Merch
	query := "SELECT id, name, price FROM merch"
	err := r.db.Select(&merchs, query)
	return merchs, err
}

func (r *merchRepository) GetByID(id int) (*models.Merch, error) {
	var merch models.Merch
	query := "SELECT id, name, price FROM merch WHERE id=$1"
	err := r.db.Get(&merch, query, id)
	if err != nil {
		return nil, err
	}
	return &merch, nil
}

func (r *merchRepository) GetByName(name string) (*models.Merch, error) {
	var merch models.Merch
	query := "SELECT id, name, price FROM merch WHERE name=$1"
	err := r.db.Get(&merch, query, name)
	if err != nil {
		return nil, err
	}
	return &merch, nil
}
