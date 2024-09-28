package repository

import (
	"testing"

	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAccountRepository_Create(t *testing.T) {

	ResetTestDB()

	repo := NewAccountRepository(db)

	account := &model.Account{
		DocumentNumber: "123456780",
	}

	accountError := &model.Account{
		DocumentNumber: "123456780",
	}

	createdAccount, err := repo.Create(account)

	assert.NoError(t, err)
	assert.NotZero(t, createdAccount.ID)
	assert.Equal(t, "123456780", createdAccount.DocumentNumber)

	_, err = repo.Create(accountError)
	assert.Error(t, err)

}

func TestAccountRepository_FindById(t *testing.T) {

	ResetTestDB()

	repo := NewAccountRepository(db)

	account := &model.Account{
		DocumentNumber: "123456",
	}

	createdAccount, err := repo.Create(account)
	assert.NoError(t, err)

	foundAccount, err := repo.FindById(createdAccount.ID)
	assert.NoError(t, err)

	assert.Equal(t, createdAccount.ID, foundAccount.ID)
	assert.Equal(t, "123456", foundAccount.DocumentNumber)
}

func TestAccountRepository_FindById_NotFound(t *testing.T) {

	ResetTestDB()

	repo := NewAccountRepository(db)

	_, err := repo.FindById(999)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
