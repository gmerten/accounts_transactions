package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dto "github.com/gmerten/accounts_transactions/api/dto"
	"github.com/gmerten/accounts_transactions/api/mapper"
	internalErrors "github.com/gmerten/accounts_transactions/internal/error"
	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionHandler_CreateTransactionSuccess(t *testing.T) {
	mockTransactionService := new(MockTransactionService)
	mockAccountService := new(MockAccountService)
	handler := NewTransactionHandler(mockTransactionService, mockAccountService)

	createTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1,
		OperationTypeID: 1,
	}

	account := &model.Account{
		ID:             1,
		DocumentNumber: "12345678",
	}

	transaction := &model.Transaction{
		ID:              1,
		AccountID:       1,
		Amount:          1,
		TransactionDate: time.Now(),
		OperationType:   model.Purchase,
	}

	mockAccountService.On("GetAccountById", int64(1)).Return(account, nil)

	mockTransactionService.On("CreateTransaction", mock.AnythingOfType("*model.Transaction")).Return(transaction, nil)

	createTransactionJSON, _ := json.Marshal(createTransactionRequest)

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(createTransactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateTransaction(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var returnedTransaction dto.CreateTransactionResponse
	_ = json.NewDecoder(rr.Body).Decode(&returnedTransaction)

	assert.Equal(t, account.ID, returnedTransaction.AccountID)
	assert.Equal(t, transaction.Amount, returnedTransaction.Amount)

	mockTransactionService.AssertExpectations(t)
	mockAccountService.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransactionInvalidJSONError(t *testing.T) {
	mockTransactionService := new(MockTransactionService)
	mockAccountService := new(MockAccountService)
	handler := NewTransactionHandler(mockTransactionService, mockAccountService)

	createTransactionJSON, _ := json.Marshal("invalid json")

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(createTransactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateTransaction(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestTransactionHandler_CreateTransactionInvalidOperationTypeError(t *testing.T) {
	mockTransactionService := new(MockTransactionService)
	mockAccountService := new(MockAccountService)
	handler := NewTransactionHandler(mockTransactionService, mockAccountService)

	createTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1,
		OperationTypeID: 6,
	}

	createTransactionJSON, _ := json.Marshal(createTransactionRequest)

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(createTransactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateTransaction(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockTransactionService.AssertExpectations(t)
	mockAccountService.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransactionAccountCustomError(t *testing.T) {
	mockTransactionService := new(MockTransactionService)
	mockAccountService := new(MockAccountService)
	handler := NewTransactionHandler(mockTransactionService, mockAccountService)

	createTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1,
		OperationTypeID: 1,
	}

	mockAccountService.On("GetAccountById", int64(1)).Return(nil, internalErrors.NewNotFoundError("account not found"))

	createTransactionJSON, _ := json.Marshal(createTransactionRequest)

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(createTransactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateTransaction(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	mockAccountService.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransactionAccountGenericError(t *testing.T) {
	mockTransactionService := new(MockTransactionService)
	mockAccountService := new(MockAccountService)
	handler := NewTransactionHandler(mockTransactionService, mockAccountService)

	createTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1,
		OperationTypeID: 1,
	}

	mockAccountService.On("GetAccountById", int64(1)).Return(nil, errors.New("generic error"))

	createTransactionJSON, _ := json.Marshal(createTransactionRequest)

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(createTransactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateTransaction(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockAccountService.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransactionError(t *testing.T) {
	mockTransactionService := new(MockTransactionService)
	mockAccountService := new(MockAccountService)
	handler := NewTransactionHandler(mockTransactionService, mockAccountService)

	createTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1,
		OperationTypeID: 1,
	}

	account := &model.Account{
		ID:             1,
		DocumentNumber: "12345678",
	}

	mockAccountService.On("GetAccountById", int64(1)).Return(account, nil)

	mockTransactionService.On("CreateTransaction", mock.AnythingOfType("*model.Transaction")).Return(nil, errors.New("error creating account"))

	createTransactionJSON, _ := json.Marshal(createTransactionRequest)

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(createTransactionJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateTransaction(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockTransactionService.AssertExpectations(t)
	mockAccountService.AssertExpectations(t)
}

func TestTransactionHandler_TransactionMapper(t *testing.T) {

	createPurchaseTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1000,
		OperationTypeID: 1,
	}

	createPaymentTransactionRequest := dto.CreateTransactionRequest{
		AccountID:       1,
		Amount:          1000,
		OperationTypeID: 4,
	}

	paymentTransaction := mapper.ToTransaction(createPaymentTransactionRequest)
	purchaseTransaction := mapper.ToTransaction(createPurchaseTransactionRequest)

	assert.Equal(t, -createPurchaseTransactionRequest.Amount, purchaseTransaction.Amount)
	assert.Equal(t, createPaymentTransactionRequest.Amount, paymentTransaction.Amount)

}
