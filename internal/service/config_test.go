package service

import (
	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(account *model.Account) (*model.Account, error) {
	args := m.Called(account)
	res := args.Get(0)
	err := args.Error(1)

	if err != nil {
		return nil, err
	}
	return res.(*model.Account), err
}

func (m *MockAccountRepository) FindById(accountID int64) (*model.Account, error) {
	args := m.Called(accountID)
	res := args.Get(0)
	err := args.Error(1)

	if err != nil {
		return nil, err
	}

	return res.(*model.Account), err
}

func (m *MockTransactionRepository) Create(transaction *model.Transaction) (*model.Transaction, error) {
	args := m.Called(transaction)

	res := args.Get(0)
	err := args.Error(1)

	if err != nil {
		return nil, err
	}
	return res.(*model.Transaction), err
}
