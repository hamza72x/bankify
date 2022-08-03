## Running

-   Run Docker
-   Copy `sample.env` to `app.env` and make necessary changes
-   `$ make up`
-   For testing: `make test`

## Endpoints

-   Create account

```sh
curl -L -d '{"name": "Hamza"}' localhost:3001/accounts
```

-   List accounts

```sh
curl -L localhost:3001/accounts
```

- Deposit money

```sh
curl -L -d '{"amount": 100}' localhost:3001/accounts/1/deposit
```

- Withdraw money

```sh
curl -L -d '{"amount": 100}' localhost:3001/accounts/1/withdraw
```

- Transfer money

```sh
curl -L -d '{"amount": 100, "from_account_id": 1, "to_account_id": 2}' localhost:3001/transfer
```

## Testing
