package main

import (
	"database/sql"
	"flag"
	server "hamza72x/bankify/api"
	"hamza72x/bankify/config"
	"hamza72x/bankify/db"
	sqlc "hamza72x/bankify/db/sqlc"
	"hamza72x/bankify/util"

	_ "github.com/lib/pq"
)

type flagArgs struct {
	Env string
}

var (
	flagArg flagArgs
)

func main() {

	parseFlags()

	// config load
	cfg, err := config.LoadConfig(flagArg.Env)
	util.LogFatal(err, "failed to load config")

	util.PrettyPrint(&cfg)

	// set database connection
	conn, err := sql.Open(config.DB_DRIVER, cfg.GetDBUrl())
	util.LogFatal(err, "failed to open db connection")

	// run migration
	db.RunMigration("db/migration", cfg)

	// new store
	store := sqlc.NewStore(conn)

	// new server
	server, err := server.New(cfg, store)
	util.LogFatal(err, "failed to create server")

	// start server
	err = server.Start()
	util.LogFatal(err, "failed to start server")
}

func parseFlags() {
	flag.StringVar(&flagArg.Env, "env", "app.env", "the config file name")
	flag.Parse()
}
