# Account & Transactions API

This application aims to manage accounts and transactions. It allows the creation and retrieval of accounts, as well as the ability to create transactions associated with these accounts.

## Technologies Used

- **Language**: Go
- **Database**: MySQL

## Features

- **Create Accounts**: Allows the creation of new accounts.
- **Retrieve Accounts**: Enables searching for accounts by ID.
- **Create Transactions**: Allows the creation of transactions associated with an existing account.

## How to Run

### Requirements 

- Go 1.21 (minimum version)
- Docker and Docker Compose


### Running via Docker

1. Clone this repository:

    ```bash
    git clone https://github.com/gmerten/accounts_transactions.git
    ```

2. Run the following command to start the application services:

    ```bash
    docker-compose up --build
    ```

This will automatically set up the application and the MySQL database.

### Running via IDE

If you prefer to run the application directly in your IDE, follow these steps:

1. Comment out the `app` section in the `docker-compose.yml` file and run the command `docker-compose up`

```bash
   
services:
  db:
    image: mysql:8.0.39
    container_name: mysql_local
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: transactions
      MYSQL_USER: user
      MYSQL_PASSWORD: password
    healthcheck:
      test: mysqladmin ping -h 127.0.0.1 -u $$MYSQL_USER --password=$$MYSQL_PASSWORD
      start_period: 5s
      interval: 5s
      timeout: 5s
      retries: 10
    ports:
      - "3306:3306"
    volumes:
      - ./scripts/migrations:/docker-entrypoint-initdb.d

  # app:
  #  build:
  #    context: .
  #    dockerfile: Dockerfile
  #  depends_on:
  #     db:
  #       condition: service_healthy
  #  environment:
  #    DB_HOST: db
  #    DB_PORT: 3306
  #    DB_USER: user
  #    DB_PASSWORD: password
  #    DB_NAME: transactions
  #  ports:
  #    - "8080:8080"
  #  links:
  #    - db
```

2. Configure the environment variables related to the MySQL database:

    ```bash
    DB_HOST=localhost
    DB_USER=user
    DB_PASSWORD=password
    DB_NAME=transactions
    DB_PORT=3306
    ```
   
3. Start the `./cmd/api/main.go` file from your IDE.



## API Examples

### 1. Create Account

To create a new account, use the following `curl` command:

```bash
curl --request POST \
  --url http://localhost:8080/accounts \
  --header 'Content-Type: application/json' \
  --data '{
	"document_number": "1314324234"
}'
```

### 2. Get a Account

To retrieve an account by ID, use the following curl command replacing `{accountID}` with the account ID:

```bash
curl --request GET \
  --url http://localhost:8080/accounts/{accountID}
```

### 3. Create a Transaction

To create a new transaction, use the following `curl` command:

```bash
curl --request POST \
  --url http://localhost:8080/transactions \
  --header 'Content-Type: application/json' \
  --data '{
	"account_id": 1,
	"amount": 200.50,
	"operation_type_id": 1
}'
```

## Swagger Documentation

The API has OpenAPI documentation available via Swagger, which can be accessed at:
http://localhost:8080/swagger/index.html


