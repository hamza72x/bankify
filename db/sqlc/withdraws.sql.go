// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: withdraws.sql

package db

import (
	"context"
)

const createWithdraw = `-- name: CreateWithdraw :one
INSERT INTO "withdraws" (
    "account_id", "amount"
) VALUES (
    $1, $2
)
RETURNING id, amount, account_id, created_at
`

type CreateWithdrawParams struct {
	AccountID int64   `json:"account_id"`
	Amount    float64 `json:"amount"`
}

func (q *Queries) CreateWithdraw(ctx context.Context, arg CreateWithdrawParams) (Withdraw, error) {
	row := q.db.QueryRowContext(ctx, createWithdraw, arg.AccountID, arg.Amount)
	var i Withdraw
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.AccountID,
		&i.CreatedAt,
	)
	return i, err
}
