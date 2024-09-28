package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	api "github.com/gmerten/accounts_transactions/api/dto"
	"github.com/gmerten/accounts_transactions/api/mapper"
	internalErrors "github.com/gmerten/accounts_transactions/internal/error"
	"github.com/gmerten/accounts_transactions/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type accountHandler struct {
	accountService service.AccountService
}

type AccountHandler interface {
	HandleGetAccount(w http.ResponseWriter, r *http.Request)
	HandleCreateAccount(w http.ResponseWriter, r *http.Request)
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{accountService}
}

// HandleGetAccount
// @Summary Get a account by id
// @Description This endpoint get a account by id
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountID path uint true "Account ID"
// @Success 200 {object} api.GetAccountResponse
// @Router /accounts/{accountID} [get]
func (a *accountHandler) HandleGetAccount(w http.ResponseWriter, r *http.Request) {
	accountIDParam := chi.URLParam(r, "accountID")
	accountID, err := strconv.ParseInt(accountIDParam, 10, 64)
	if err != nil {
		HandleError(w, internalErrors.NewValidationError("Invalid Operation ID"))
		return
	}

	account, err := a.accountService.GetAccountById(accountID)
	if err != nil {
		log.WithField("accountID", accountID).WithError(err).Error("Error getting account")
		_, ok := err.(CustomError)
		if ok {
			HandleError(w, err)
			return
		}
		HandleError(w, internalErrors.NewUnknownError("Error getting account"))
		return
	}

	response := mapper.ToGetAccountResponse(account)

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}

// HandleCreateAccount
// @Summary Creates a new account
// @Description This endpoint creates a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body api.CreateAccountRequest true "Request body"
// @Success 200 {object} api.CreateAccountResponse
// @Router /accounts [post]
func (a *accountHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody api.CreateAccountRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.WithError(err).Error("Error parsing request body")
		HandleError(w, internalErrors.NewValidationError("Invalid request body"))
		return
	}

	err = validator.New().Struct(requestBody)
	if err != nil {
		log.WithError(err).Error("Error validating request body")
		HandleError(w, internalErrors.NewValidationError("Invalid request body"))
		return
	}

	account := mapper.ToAccount(requestBody)

	account, err = a.accountService.CreateAccount(account)

	if err != nil {
		log.WithError(err).Error("Error creating account")
		_, ok := err.(CustomError)
		if ok {
			HandleError(w, err)
			return
		}
		HandleError(w, internalErrors.NewUnknownError("Error creating account"))
		return
	}

	response := mapper.ToCreateAccountResponse(account)

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(response)
}
