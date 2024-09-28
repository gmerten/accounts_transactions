package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/gmerten/accounts_transactions/api/dto"
	internalErrors "github.com/gmerten/accounts_transactions/internal/error"
	"github.com/gmerten/accounts_transactions/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestAccountHandler_CreateAccountSuccess(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	createAccountRequest := dto.CreateAccountRequest{
		DocumentNumber: "12345678",
	}

	account := &model.Account{
		DocumentNumber: "12345678",
	}

	mockService.On("CreateAccount", account).Return(account, nil)

	createAccountJSON, _ := json.Marshal(createAccountRequest)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(createAccountJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateAccount(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestAccountHandler_CreateAccountInvalidJSONError(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	createAccountJSON, _ := json.Marshal("invalid json")

	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(createAccountJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateAccount(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAccountHandler_CreateAccountError(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	createAccountRequest := dto.CreateAccountRequest{
		DocumentNumber: "12345678",
	}

	account := &model.Account{
		DocumentNumber: "12345678",
	}

	mockService.On("CreateAccount", account).Return(nil, errors.New("error creating account"))

	createAccountJSON, _ := json.Marshal(createAccountRequest)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(createAccountJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateAccount(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockService.AssertExpectations(t)
}

func TestAccountHandler_CreateAccountInvalidRequest(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	createAccountRequest := dto.CreateAccountRequest{
		DocumentNumber: "",
	}

	createAccountJSON, _ := json.Marshal(createAccountRequest)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(createAccountJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateAccount(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertExpectations(t)
}

func TestAccountHandler_CreateAccountCustomError(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	createAccountRequest := dto.CreateAccountRequest{
		DocumentNumber: "12345678",
	}

	account := &model.Account{
		DocumentNumber: "12345678",
	}

	mockService.On("CreateAccount", account).Return(nil, internalErrors.NewConflictError("account already exists"))

	createAccountJSON, _ := json.Marshal(createAccountRequest)

	req, err := http.NewRequest("POST", "/accounts", bytes.NewBuffer(createAccountJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleCreateAccount(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	mockService.AssertExpectations(t)
}

func TestAccountHandler_GetAccountSuccess(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	account := &model.Account{
		ID:             1,
		DocumentNumber: "12345678",
	}

	mockService.On("GetAccountById", int64(1)).Return(account, nil)

	req, err := http.NewRequest("GET", "/accounts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("accountID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleGetAccount(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var returnedAccount dto.GetAccountResponse
	_ = json.NewDecoder(rr.Body).Decode(&returnedAccount)

	assert.Equal(t, account.ID, returnedAccount.ID)

	mockService.AssertExpectations(t)
}

func TestAccountHandler_GetAccountInvalidIDError(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	req, err := http.NewRequest("GET", "/accounts/aaaa", nil)
	if err != nil {
		t.Fatal(err)
	}

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("accountID", "aaa")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleGetAccount(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockService.AssertExpectations(t)
}

func TestAccountHandler_GetAccountCustomError(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	mockService.On("GetAccountById", int64(1)).Return(nil, internalErrors.NewNotFoundError("account not found"))

	req, err := http.NewRequest("GET", "/accounts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("accountID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleGetAccount(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	mockService.AssertExpectations(t)
}

func TestAccountHandler_GetAccountInternalServerError(t *testing.T) {
	mockService := new(MockAccountService)
	handler := NewAccountHandler(mockService)

	mockService.On("GetAccountById", int64(1)).Return(nil, errors.New("generic error"))

	req, err := http.NewRequest("GET", "/accounts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("accountID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleGetAccount(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}
