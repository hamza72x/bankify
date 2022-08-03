-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccounts :many
SELECT * FROM accounts;

-- name: CreateAccount :one
INSERT INTO accounts (
  name, balance
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateBalance :one
UPDATE "accounts" SET
"balance" = $1
WHERE "id" = $2
RETURNING *;
