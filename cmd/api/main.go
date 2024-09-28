package main

import (
	"net/http"

	api "github.com/gmerten/accounts_transactions/api/handler"
	_ "github.com/gmerten/accounts_transactions/docs"
	"github.com/gmerten/accounts_transactions/internal/config"
	"github.com/gmerten/accounts_transactions/internal/repository"
	"github.com/gmerten/accounts_transactions/internal/service"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Accounts & Transactions API
// @version 1.0
// @description API to manage accounts and transactions
// @host localhost:8080
// @BasePath /
func main() {

	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)

	router := chi.NewRouter()

	db := config.GetDBConnection()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)
	accountHandler := api.NewAccountHandler(accountService)

	transactionRepository := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepository)
	transactionHandler := api.NewTransactionHandler(transactionService, accountService)

	router.Get("/accounts/{accountID}", accountHandler.HandleGetAccount)
	router.Post("/accounts", accountHandler.HandleCreateAccount)
	router.Post("/transactions", transactionHandler.HandleCreateTransaction)
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8080", router))

}
