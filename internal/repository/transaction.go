package repository

import (
	"github.com/gmerten/accounts_transactions/internal/model"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

type TransactionRepository interface {
	Create(transaction *model.Transaction) (*model.Transaction, error)
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(transaction *model.Transaction) (*model.Transaction, error) {
	if err := r.db.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}
