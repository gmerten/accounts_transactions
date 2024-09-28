package service

import (
	"errors"
	"testing"

	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)

	transaction := &model.Transaction{
		ID:            1,
		AccountID:     1,
		Amount:        100.0,
		OperationType: 1,
	}

	mockRepo.On("Create", transaction).Return(transaction, nil)

	service := NewTransactionService(mockRepo)

	createdTransaction, err := service.CreateTransaction(transaction)

	assert.NoError(t, err)
	assert.Equal(t, transaction, createdTransaction)

	mockRepo.AssertExpectations(t)
}

func TestTransactionService_CreateTransactionError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)

	transaction := &model.Transaction{
		ID:            1,
		AccountID:     1,
		Amount:        100.0,
		OperationType: 1,
	}

	mockRepo.On("Create", transaction).Return(nil, errors.New("error creating transaction"))

	service := NewTransactionService(mockRepo)

	_, err := service.CreateTransaction(transaction)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
