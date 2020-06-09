package postgres

import (
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/prometheus/common/log"
)

const migrationsSource = "file:///Users/vladislav_tokarev/go/src/github.com/odahu/odahu-flow/packages/operator/pkg/database/migrations/postgres"

type Migrator struct {
	ConnString string
}

func (m Migrator) MigrateToLatest() error {

	mClient, err := migrate.New(migrationsSource, m.ConnString)
	if err != nil {
		return err
	}

	version, _, _ := mClient.Version()
	err = mClient.Up()
	newVersion, _, _ := mClient.Version()
	if err != nil {
		switch {
		case err == migrate.ErrNoChange || err == migrate.ErrNilVersion:
			log.Infof(
				"No migrations was applied. Old version: %d, new version: %d",
				version, newVersion,
			)
		default:
			return err
		}
	}

	return nil
}

func (m Migrator) DownAllMigrations() error {

	mClient, err := migrate.New(migrationsSource, m.ConnString)
	if err != nil {
		return err
	}

	err = mClient.Down()
	if err != nil {
		return err
	}

	return nil
}
