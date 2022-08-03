package api

import (
	"database/sql"
	"errors"
	"hamza72x/bankify/api/reqs"
	db "hamza72x/bankify/db/sqlc"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerAccountRoutes(r *gin.Engine) {

	accounts := r.Group("/accounts")

	accounts.POST("/", s.createAccount)
	accounts.GET("/", s.getAccounts)
	accounts.POST("/:id/deposit", s.deposit)
	accounts.POST("/:id/withdraw", s.withdraw)
}

func (s *Server) createAccount(c *gin.Context) {
	var req reqs.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	acc, err := s.store.CreateAccount(c, db.CreateAccountParams{Name: req.Name, Balance: 0})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, acc)
}

// get accounts
// TODO :- must have admin auth, but for now - it's okay
func (s *Server) getAccounts(c *gin.Context) {

	list, err := s.store.GetAccounts(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, list)
}

// deposit
func (s *Server) deposit(c *gin.Context) {

	var reqBody reqs.WithdrawRequestBody
	var reqUri reqs.WithdrawRequestUri

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, valid := s.validateAccount(c, reqUri.AccountId)

	if !valid {
		return
	}

	// TODO :- implement auth

	err := s.store.ExecTx(c, func(txQuery *db.Queries) error {

		reqBody.Amount = math.Abs(reqBody.Amount)

		// create deposit history
		_, err := txQuery.CreateDeposit(c, db.CreateDepositParams{
			AccountID: reqUri.AccountId,
			Amount:    reqBody.Amount,
		})

		if err != nil {
			return err
		}

		// remove amount from "fromAccount"
		_, err = txQuery.UpdateBalance(c, db.UpdateBalanceParams{
			Balance: account.Balance + reqBody.Amount,
			ID:      account.ID,
		})

		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	account, err = s.store.GetAccount(c, reqUri.AccountId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deposit successful",
		"account": account,
	})
}

// withdraw
func (s *Server) withdraw(c *gin.Context) {
	var reqBody reqs.WithdrawRequestBody
	var reqUri reqs.WithdrawRequestUri

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, valid := s.validateAccount(c, reqUri.AccountId)

	if !valid {
		return
	}

	// TODO :- implement auth

	if reqBody.Amount > account.Balance {
		c.JSON(http.StatusBadRequest, errResponse(errors.New("insufficient balance")))
		return
	}

	err := s.store.ExecTx(c, func(txQuery *db.Queries) error {

		reqBody.Amount = math.Abs(reqBody.Amount)

		// create withdraw history
		_, err := txQuery.CreateWithdraw(c, db.CreateWithdrawParams{
			AccountID: reqUri.AccountId,
			Amount:    reqBody.Amount,
		})

		if err != nil {
			return err
		}

		// remove amount from "fromAccount"
		_, err = txQuery.UpdateBalance(c, db.UpdateBalanceParams{
			Balance: account.Balance - reqBody.Amount,
			ID:      account.ID,
		})

		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	account, err = s.store.GetAccount(c, reqUri.AccountId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "withdraw successful",
		"account": account,
	})
}

func (s *Server) validateAccount(c *gin.Context, accountID int64) (db.Account, bool) {
	account, err := s.store.GetAccount(c, accountID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}

	return account, true
}
