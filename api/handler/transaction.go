package api

import (
	"encoding/json"
	"net/http"

	"github.com/gmerten/accounts_transactions/api/dto"
	"github.com/gmerten/accounts_transactions/api/mapper"
	internalErrors "github.com/gmerten/accounts_transactions/internal/error"
	"github.com/gmerten/accounts_transactions/internal/service"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type transactionHandler struct {
	transactionService service.TransactionService
	accountService     service.AccountService
}

type TransactionHandler interface {
	HandleCreateTransaction(w http.ResponseWriter, r *http.Request)
}

func NewTransactionHandler(transactionService service.TransactionService, accountService service.AccountService) TransactionHandler {
	return &transactionHandler{transactionService,
		accountService}
}

// HandleCreateTransaction
// @Summary Creates a new transaction
// @Description This endpoint creates a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body api.CreateTransactionRequest true "Request body"
// @Success 200 {object} api.CreateTransactionResponse
// @Router /transactions [post]
func (t *transactionHandler) HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody api.CreateTransactionRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.WithError(err).Error("Error decoding request body")
		HandleError(w, internalErrors.NewValidationError("Invalid Request Body"))
		return
	}

	err = validator.New().Struct(requestBody)
	if err != nil {
		log.WithError(err).Error("Error validating request body")
		HandleError(w, internalErrors.NewValidationError("Invalid request body"))
		return
	}

	transaction := mapper.ToTransaction(requestBody)

	_, err = t.accountService.GetAccountById(transaction.AccountID)

	if err != nil {
		log.WithError(err).Error("Error getting account")
		_, ok := err.(CustomError)
		if ok {
			HandleError(w, err)
			return
		}
		HandleError(w, internalErrors.NewUnknownError("Error getting account"))
		return
	}

	transaction, err = t.transactionService.CreateTransaction(transaction)

	if err != nil {
		log.WithField("accountID", requestBody.AccountID).WithError(err).Error("Error creating transaction")
		HandleError(w, internalErrors.NewUnknownError("Fail creating transaction"))
		return
	}

	response := mapper.ToCreateTransactionResponse(transaction)

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(response)

}
