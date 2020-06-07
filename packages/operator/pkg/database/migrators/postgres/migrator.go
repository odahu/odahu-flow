package postgres

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type Migrator struct {
	ConnString string
}

func (m Migrator) MigrateToLatest() error {

	// TODO: corner cases

	fmt.Print("Migrate to latest")
	migrateSource := "file:///Users/vladislav_tokarev/go/src/github.com/odahu/odahu-flow/packages/operator/pkg/database/migrations/postgres"
	mClient, err := migrate.New(migrateSource, m.ConnString)
	if err != nil {
		return err
	}

	err = mClient.Up()
	if err != nil {
		return err
	}

	return nil
}
