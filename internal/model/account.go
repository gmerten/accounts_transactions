package model

type Account struct {
	ID             int64         `gorm:"primaryKey"`
	DocumentNumber string        `gorm:"unique;not null"`
	Transactions   []Transaction `gorm:"foreignKey:AccountID;references:ID"`
}
