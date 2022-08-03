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

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {

	server := newTestServer(t)
	accName := "Hamza_TestCreateAccount"

	testCases := []struct {
		name          string
		buildRequest  func() (*db.Account, *http.Request)
		checkResponse func(t *testing.T, rr *httptest.ResponseRecorder) []db.Account
	}{
		{
			name: "Successful Creation",
			buildRequest: func() (*db.Account, *http.Request) {
				return nil, httptest.NewRequest(
					http.MethodPost,
					"/accounts/",
					bytes.NewBufferString(fmt.Sprintf(`{"name":"%s"}`, accName)),
				)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) []db.Account {
				var acc db.Account
				err := json.Unmarshal(rr.Body.Bytes(), &acc)
				require.NoError(t, err)

				require.Greater(t, acc.ID, int64(0))
				require.Equal(t, http.StatusOK, rr.Code)
				require.Equal(t, accName, acc.Name)
				require.Equal(t, float64(0), acc.Balance)

				return []db.Account{acc}
			},
		},
		{
			name: "Uniqueness Check",
			buildRequest: func() (*db.Account, *http.Request) {
				acc, err := server.store.CreateAccount(
					context.Background(),
					db.CreateAccountParams{Name: accName, Balance: 100},
				)

				require.NoError(t, err)

				return &acc, httptest.NewRequest(
					http.MethodPost,
					"/accounts/",
					bytes.NewBufferString(fmt.Sprintf(`{"name":"%s"}`, accName)),
				)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) []db.Account {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
				return nil
			},
		},
		{
			name: "Missing Name",
			buildRequest: func() (*db.Account, *http.Request) {
				return nil, httptest.NewRequest(
					http.MethodPost,
					"/accounts/",
					bytes.NewBufferString(`{}`),
				)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) []db.Account {
				require.Equal(t, http.StatusBadRequest, rr.Code)
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			acc, req := tc.buildRequest()
			server.router.ServeHTTP(rr, req)

			accs := tc.checkResponse(t, rr)

			if acc != nil {
				tearDown(argTearDown{
					t:        t,
					store:    server.store,
					table:    "accounts",
					column:   "id",
					relation: "=",
					value:    acc.ID,
				})

			}
			for _, acc := range accs {
				tearDown(argTearDown{
					t:        t,
					store:    server.store,
					table:    "accounts",
					column:   "id",
					relation: "=",
					value:    acc.ID,
				})
			}
		})
	}
}
