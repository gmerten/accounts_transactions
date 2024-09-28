package mapper

import (
	"time"

	api "github.com/gmerten/accounts_transactions/api/dto"
	"github.com/gmerten/accounts_transactions/internal/model"
)

func ToCreateAccountResponse(account *model.Account) api.CreateAccountResponse {
	return api.CreateAccountResponse{
		DocumentNumber: account.DocumentNumber,
		ID:             account.ID,
	}
}

func ToGetAccountResponse(account *model.Account) api.GetAccountResponse {
	return api.GetAccountResponse{
		DocumentNumber: account.DocumentNumber,
		ID:             account.ID,
	}
}

func ToCreateTransactionResponse(transaction *model.Transaction) api.CreateTransactionResponse {
	return api.CreateTransactionResponse{
		TransactionID:   transaction.ID,
		AccountID:       transaction.AccountID,
		Amount:          transaction.Amount,
		OperationTypeID: uint(transaction.OperationType),
	}
}

func ToAccount(request api.CreateAccountRequest) *model.Account {
	return &model.Account{
		DocumentNumber: request.DocumentNumber,
	}
}

func ToTransaction(request api.CreateTransactionRequest) *model.Transaction {

	operationType := model.OperationType(request.OperationTypeID)
	amount := normalizeAmount(operationType, request.Amount)

	return &model.Transaction{
		OperationType:   operationType,
		Amount:          amount,
		AccountID:       request.AccountID,
		TransactionDate: time.Now(),
	}
}

func normalizeAmount(operationType model.OperationType, amount float64) float64 {
	switch operationType {
	case model.Payment:
		return amount
	default:
		return -amount
	}
}
