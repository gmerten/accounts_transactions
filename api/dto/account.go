package api

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" validate:"required"`
}

type CreateAccountResponse struct {
	DocumentNumber string `json:"document_number"`
	ID             int64  `json:"account_id"`
}

type GetAccountResponse struct {
	DocumentNumber string `json:"document_number"`
	ID             int64  `json:"account_id"`
}
