package model

import (
	"time"
)

type OperationType int

const (
	Purchase OperationType = iota + 1
	InstallmentPurchase
	Withdrawal
	Payment
)

type Transaction struct {
	ID              int64 `gorm:"primaryKey"`
	OperationType   OperationType
	Amount          float64
	TransactionDate time.Time
	AccountID       int64
	Account         Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
