package reqs

type TransferRequest struct {
	FromAccountId int64   `json:"from_account_id" binding:"required"`
	ToAccountId   int64   `json:"to_account_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}
