package api

import (
	"errors"
	"hamza72x/bankify/api/reqs"
	db "hamza72x/bankify/db/sqlc"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerTransferRoutes(r *gin.Engine) {
	r.POST("/transfer", s.transfer)
}

func (s *Server) transfer(c *gin.Context) {
	var req reqs.TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	fromAccount, valid := s.validateAccount(c, req.FromAccountId)

	if !valid {
		return
	}

	toAccount, valid := s.validateAccount(c, req.ToAccountId)

	if !valid {
		return
	}

	// TODO :- implement auth for "fromAccount"

	if req.Amount > fromAccount.Balance {
		c.JSON(http.StatusBadRequest, errResponse(errors.New("insufficient balance")))
		return
	}

	var transfer db.Transfer

	err := s.store.ExecTx(c, func(txQuery *db.Queries) error {

		var err error
		req.Amount = math.Abs(req.Amount)

		// create transfer history
		transfer, err = txQuery.CreateTransfer(c, db.CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        req.Amount,
		})

		if err != nil {
			return err
		}

		// remove amount from "fromAccount"
		_, err = txQuery.UpdateBalance(c, db.UpdateBalanceParams{
			Balance: fromAccount.Balance - req.Amount,
			ID:      fromAccount.ID,
		})

		if err != nil {
			return err
		}

		// add amount to "toAccount"
		_, err = txQuery.UpdateBalance(c, db.UpdateBalanceParams{
			Balance: toAccount.Balance + req.Amount,
			ID:      toAccount.ID,
		})

		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "transfer successful",
		"transfer": transfer,
	})
}
