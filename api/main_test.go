package api

import (
	"database/sql"
	"fmt"
	"hamza72x/bankify/config"
	"hamza72x/bankify/db"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	sqlc "hamza72x/bankify/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *Server {
	// config load
	cfg, err := config.LoadConfig("../app.env")
	require.NoError(t, err)

	// set database connection
	conn, err := sql.Open(config.DB_DRIVER, cfg.GetDBUrl())
	require.NoError(t, err)

	// run migration
	db.RunMigration("../db/migration", cfg)

	// new store
	store := sqlc.NewStore(conn)

	// new server
	server, err := New(cfg, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func execRequest(t *testing.T, server *Server, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)
	return rr
}

type argTearDown struct {
	t        *testing.T
	store    *sqlc.Store
	table    string
	column   string
	relation string
	value    interface{}
}

func tearDown(arg argTearDown) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s %s $1", arg.table, arg.column, arg.relation)

	_, err := arg.store.DBExecRaw(query, arg.value)

	require.NoError(arg.t, err)
}
