-- name: CreateDeposit :one
INSERT INTO "deposits" (
    "account_id", "amount"
) VALUES (
    $1, $2
)
RETURNING *;
