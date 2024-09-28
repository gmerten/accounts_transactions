package service

import (
	"errors"
	"testing"

	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAccountService_CreateAccount(t *testing.T) {

	mockRepo := new(MockAccountRepository)

	account := &model.Account{
		ID:             1,
		DocumentNumber: "12345678",
	}

	mockRepo.On("Create", account).Return(account, nil)

	service := NewAccountService(mockRepo)
	createdAccount, err := service.CreateAccount(account)

	assert.NoError(t, err)
	assert.Equal(t, account, createdAccount)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_CreateAccountGenericError(t *testing.T) {

	mockRepo := new(MockAccountRepository)

	account := &model.Account{
		ID:             2,
		DocumentNumber: "12345678",
	}

	mockRepo.On("Create", account).Return(nil, errors.New("error creating account"))

	service := NewAccountService(mockRepo)

	_, err := service.CreateAccount(account)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_CreateAccountMySQLDuplicatedKeyError(t *testing.T) {

	mockRepo := new(MockAccountRepository)

	account := &model.Account{
		ID:             2,
		DocumentNumber: "12345678",
	}

	duplicatedKeyError := &mysql.MySQLError{
		Number: 1062,
	}

	mockRepo.On("Create", account).Return(nil, duplicatedKeyError)

	service := NewAccountService(mockRepo)

	_, err := service.CreateAccount(account)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_CreateAccountDBDuplicatedKeyError(t *testing.T) {

	mockRepo := new(MockAccountRepository)

	account := &model.Account{
		ID:             2,
		DocumentNumber: "12345678",
	}

	mockRepo.On("Create", account).Return(nil, gorm.ErrDuplicatedKey)

	service := NewAccountService(mockRepo)

	_, err := service.CreateAccount(account)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_GetAccountById(t *testing.T) {

	mockRepo := new(MockAccountRepository)

	account := &model.Account{
		ID:             1,
		DocumentNumber: "123435435",
	}
	mockRepo.On("FindById", int64(1)).Return(account, nil)

	service := NewAccountService(mockRepo)

	foundAccount, err := service.GetAccountById(1)
	assert.NoError(t, err)
	assert.Equal(t, account, foundAccount)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_GetAccountByIdGenericError(t *testing.T) {

	mockRepo := new(MockAccountRepository)

	mockRepo.On("FindById", int64(2)).Return(nil, errors.New("generic error"))

	service := NewAccountService(mockRepo)

	_, err := service.GetAccountById(2)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAccountService_GetAccountByIdNotFoundError(t *testing.T) {

	mockRepo := new(MockAccountRepository)

	mockRepo.On("FindById", int64(2)).Return(nil, gorm.ErrRecordNotFound)

	service := NewAccountService(mockRepo)

	_, err := service.GetAccountById(2)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
