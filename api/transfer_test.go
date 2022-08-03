package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	db "hamza72x/bankify/db/sqlc"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {

	server := newTestServer(t)

	fromAcc, err := server.store.CreateAccount(
		context.Background(),
		db.CreateAccountParams{Name: "Hamza_TestTransfer", Balance: 100},
	)
	require.NoError(t, err)

	toAcc, err := server.store.CreateAccount(
		context.Background(),
		db.CreateAccountParams{Name: "Ahmed_TestTransfer", Balance: 100},
	)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		buildRequest  func() *http.Request
		checkResponse func(t *testing.T, rr *httptest.ResponseRecorder) db.Transfer
	}{
		{
			name: "insufficient balance",
			buildRequest: func() *http.Request {
				return httptest.NewRequest(
					http.MethodPost,
					"/transfer",
					bytes.NewBufferString(fmt.Sprintf(
						`{"from_account_id":%d,"to_account_id":%d,"amount":200}`,
						fromAcc.ID,
						toAcc.ID,
					)),
				)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) db.Transfer {
				require.Equal(t, http.StatusBadRequest, rr.Code)

				var response gin.H

				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, "insufficient balance", response["error"])

				return db.Transfer{}
			},
		},
		{
			name: "successful transfer",
			buildRequest: func() *http.Request {
				return httptest.NewRequest(
					http.MethodPost,
					"/transfer",
					bytes.NewBufferString(fmt.Sprintf(
						`{"from_account_id":%d,"to_account_id":%d,"amount":100}`,
						fromAcc.ID,
						toAcc.ID,
					)),
				)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) db.Transfer {
				require.Equal(t, http.StatusOK, rr.Code)

				var response struct {
					Message  string      `json:"message"`
					Transfer db.Transfer `json:"transfer"`
				}

				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, "transfer successful", response.Message)
				require.Equal(t, fromAcc.ID, response.Transfer.FromAccountID)
				require.Equal(t, toAcc.ID, response.Transfer.ToAccountID)
				require.Equal(t, float64(100), response.Transfer.Amount)

				return response.Transfer
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			req := tc.buildRequest()
			server.router.ServeHTTP(rr, req)

			tc.checkResponse(t, rr)

			tearDown(argTearDown{
				t:        t,
				store:    server.store,
				table:    "transfers",
				column:   "from_account_id",
				relation: "=",
				value:    fromAcc.ID,
			})
		})
	}

	tearDown(argTearDown{
		t:        t,
		store:    server.store,
		table:    "accounts",
		column:   "id",
		relation: "=",
		value:    fromAcc.ID,
	})

	tearDown(argTearDown{
		t:        t,
		store:    server.store,
		table:    "accounts",
		column:   "id",
		relation: "=",
		value:    toAcc.ID,
	})
}
