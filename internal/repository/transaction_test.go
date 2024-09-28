package repository

import (
	"testing"
	"time"

	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_Create(t *testing.T) {

	ResetTestDB()

	accountRepo := NewAccountRepository(db)
	repo := NewTransactionRepository(db)

	account := &model.Account{
		DocumentNumber: "123456780",
	}

	createdAccount, err := accountRepo.Create(account)

	assert.NoError(t, err)
	assert.NotZero(t, createdAccount.ID)

	transaction := &model.Transaction{
		AccountID:       createdAccount.ID,
		Amount:          1000,
		TransactionDate: time.Now(),
		OperationType:   model.Purchase,
	}

	createdTransaction, err := repo.Create(transaction)

	assert.NoError(t, err)
	assert.NotZero(t, createdTransaction.ID)

}

func TestTransactionRepository_CreateError(t *testing.T) {

	ResetTestDB()

	repo := NewTransactionRepository(db)

	transaction := &model.Transaction{
		AccountID:       999,
		Amount:          1000,
		TransactionDate: time.Now(),
		OperationType:   model.Purchase,
	}

	_, err := repo.Create(transaction)
	assert.Error(t, err)
}
