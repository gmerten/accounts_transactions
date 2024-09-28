package api

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id" validate:"required,gte=1"`
	Amount          float64 `json:"amount" validate:"required,gte=0"`
	OperationTypeID uint    `json:"operation_type_id" validate:"required,oneof=1 2 3 4"`
}

type CreateTransactionResponse struct {
	TransactionID   int64   `json:"transaction_id"`
	AccountID       int64   `json:"account_id"`
	Amount          float64 `json:"amount"`
	OperationTypeID uint    `json:"operation_type_id"`
}
