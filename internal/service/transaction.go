package service

import (
	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/gmerten/accounts_transactions/internal/repository"
)

type transactionService struct {
	repository repository.TransactionRepository
}

type TransactionService interface {
	CreateTransaction(transaction *model.Transaction) (*model.Transaction, error)
}

func NewTransactionService(repository repository.TransactionRepository) TransactionService {
	return &transactionService{repository}
}

func (t *transactionService) CreateTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	return t.repository.Create(transaction)
}
