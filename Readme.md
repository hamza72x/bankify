## Running

-   Run Docker
-   Copy `sample.env` to `app.env` and make necessary changes
-   `$ make up`
-   For testing: `make test`

## Endpoints

1. Create account
```sh
curl -L -d '{"name": "Hamza"}' localhost:3001/accounts
```

2. Deposit money
```sh
curl -L -d '{"amount": 100}' localhost:3001/accounts/1/deposit
```

3. Withdraw money
```sh
curl -L -d '{"amount": 100}' localhost:3001/accounts/1/withdraw
```

4. Transfer money
```sh
curl -L -d '{"amount": 100, "from_account_id": 1, "to_account_id": 2}' localhost:3001/transfer
```

## Testing
Currently no test files are provided due to the lack of time.
Will add test files if mandatory, please let me know.
