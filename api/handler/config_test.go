package api

import (
	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

type MockTransactionService struct {
	mock.Mock
}

func (m *MockAccountService) CreateAccount(account *model.Account) (*model.Account, error) {
	args := m.Called(account)
	res := args.Get(0)
	err := args.Error(1)

	if err != nil {
		return nil, err
	}
	return res.(*model.Account), err
}

func (m *MockAccountService) GetAccountById(accountId int64) (*model.Account, error) {
	args := m.Called(accountId)
	res := args.Get(0)
	err := args.Error(1)

	if err != nil {
		return nil, err
	}

	return res.(*model.Account), err
}

func (m *MockTransactionService) CreateTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	args := m.Called(transaction)

	res := args.Get(0)
	err := args.Error(1)

	if err != nil {
		return nil, err
	}
	return res.(*model.Transaction), err
}
