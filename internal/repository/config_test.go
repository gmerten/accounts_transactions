package repository

import (
	"os"
	"testing"

	"github.com/gmerten/accounts_transactions/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupTestDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err = db.AutoMigrate(&model.Account{}, &model.Transaction{}); err != nil {
		panic("failed to migrate database")
	}
}

func ResetTestDB() {
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM transactions")
}

func TestMain(m *testing.M) {
	SetupTestDB()
	code := m.Run()
	os.Exit(code)
}
