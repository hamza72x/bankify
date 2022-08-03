package reqs

type CreateAccountRequest struct {
	Name string `json:"name" binding:"required"`
}

type WithdrawRequestUri struct {
	AccountId int64 `uri:"id" binding:"required"`
}

type DepositRequestUri struct {
	AccountId int64 `uri:"id" binding:"required"`
}

type WithdrawRequestBody struct {
	Amount float64 `json:"amount" binding:"required"`
}

type DepositRequestBody struct {
	Amount float64 `json:"amount" binding:"required"`
}
