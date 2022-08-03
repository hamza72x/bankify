-- name: CreateWithdraw :one
INSERT INTO "withdraws" (
    "account_id", "amount"
) VALUES (
    $1, $2
)
RETURNING *;
