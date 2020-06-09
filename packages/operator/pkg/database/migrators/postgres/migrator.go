package postgres

import (
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source/go_bindata"
	migrations "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrations/postgres"
	"github.com/prometheus/common/log"
)

type Migrator struct {
	ConnString string
	client     *migrate.Migrate
}

func NewMigrator(connString string) (*Migrator, error) {

	migr := &Migrator{ConnString: connString}

	// wrap assets into Resource
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)
	migr.client, err = migrate.NewWithSourceInstance("go-bindata", d, connString)

	if err != nil {
		return nil, err
	}

	return migr, nil
}

func (m Migrator) MigrateToLatest() error {

	mClient := m.client

	version, _, _ := mClient.Version()
	err := mClient.Up()
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

	mClient := m.client

	err := mClient.Down()
	if err != nil {
		return err
	}

	return nil
}
