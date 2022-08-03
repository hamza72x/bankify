package db

import (
	"hamza72x/bankify/config"
	"hamza72x/bankify/util"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	col "github.com/hamza72x/go-color"
)

// migrationPath is the path to the migration files
// usually: db/migration
func RunMigration(migrationPath string, cfg config.Config) {

	migration, err := migrate.New("file://"+migrationPath, cfg.GetDBUrl())
	util.LogFatal(err, "failed to create migration instance")

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run up", err)
	}

	log.Println(col.Green("db migrated successfully"))
}
