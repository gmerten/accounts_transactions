package repository

import (
	"github.com/gmerten/accounts_transactions/internal/model"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *model.Account) (*model.Account, error)
	FindById(id int64) (*model.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db}
}

func (r *accountRepository) Create(account *model.Account) (*model.Account, error) {
	if err := r.db.Create(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) FindById(id int64) (*model.Account, error) {
	var account model.Account
	if err := r.db.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
