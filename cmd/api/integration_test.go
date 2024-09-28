package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	dto "github.com/gmerten/accounts_transactions/api/dto"
	api "github.com/gmerten/accounts_transactions/api/handler"
	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/gmerten/accounts_transactions/internal/repository"
	"github.com/gmerten/accounts_transactions/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestE2E_CreateAccountAndTransaction(t *testing.T) {

	router := setupTest()

	createAccountRequest := dto.CreateAccountRequest{
		DocumentNumber: "12345678",
	}

	createAccountJSON, _ := json.Marshal(createAccountRequest)
	rCreateAccount := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(createAccountJSON))
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(rCreateAccount, req)

	var returnedAccount dto.CreateAccountResponse
	_ = json.NewDecoder(rCreateAccount.Body).Decode(&returnedAccount)

	assert.Equal(t, http.StatusCreated, rCreateAccount.Code)
	assert.Equal(t, createAccountRequest.DocumentNumber, returnedAccount.DocumentNumber)

	accountIDParam := strconv.FormatInt(returnedAccount.ID, 10)
	rGetAccount := httptest.NewRecorder()

	req, err = http.NewRequest("GET", "/accounts/"+accountIDParam, nil)
	if err != nil {
		t.Fatal(err)
	}

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("accountID", accountIDParam)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rGetAccount, req)

	assert.Equal(t, http.StatusOK, rGetAccount.Code)

	createTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1000,
		OperationTypeID: 4,
	}

	createTransactionJSON, _ := json.Marshal(createTransactionRequest)
	rCreateTransaction := httptest.NewRecorder()

	req, err = http.NewRequest("POST", "/transactions", bytes.NewBuffer(createTransactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rCreateTransaction, req)

	var returnedTransaction dto.CreateTransactionResponse
	_ = json.NewDecoder(rCreateTransaction.Body).Decode(&returnedTransaction)

	assert.Equal(t, returnedAccount.ID, returnedTransaction.AccountID)
	assert.Equal(t, createTransactionRequest.Amount, returnedTransaction.Amount)

}

func setupTest() *chi.Mux {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err = db.AutoMigrate(&model.Account{}, &model.Transaction{}); err != nil {
		panic("failed to migrate database")
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)
	accountHandler := api.NewAccountHandler(accountService)

	transactionRepository := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepository)
	transactionHandler := api.NewTransactionHandler(transactionService, accountService)

	router.Get("/accounts/{accountID}", accountHandler.HandleGetAccount)
	router.Post("/accounts", accountHandler.HandleCreateAccount)
	router.Post("/transactions", transactionHandler.HandleCreateTransaction)

	return router
}
