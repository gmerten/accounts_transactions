package service

import (
	"errors"

	internalErrors "github.com/gmerten/accounts_transactions/internal/error"
	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/gmerten/accounts_transactions/internal/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type accountService struct {
	repository repository.AccountRepository
}

type AccountService interface {
	CreateAccount(account *model.Account) (*model.Account, error)
	GetAccountById(accountId int64) (*model.Account, error)
}

func NewAccountService(repository repository.AccountRepository) AccountService {
	return &accountService{repository}
}

func (a *accountService) CreateAccount(account *model.Account) (*model.Account, error) {
	account, err := a.repository.Create(account)
	if err != nil {
		log.WithError(err).Error("Error saving account")
		if internalErrors.IsDuplicateKeyError(err) {
			return nil, internalErrors.NewConflictError("Account with this document number already exists")
		}

		return nil, err
	}
	return account, nil
}

func (a *accountService) GetAccountById(accountId int64) (*model.Account, error) {
	account, err := a.repository.FindById(accountId)

	if err != nil {
		log.WithError(err).Error("Error getting account")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internalErrors.NewNotFoundError("Account not found")
		}
		return nil, err
	}

	return account, nil

}
