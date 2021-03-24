package testenvs

import (
	"database/sql"
	"fmt"
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"os"
	"time"
)

// SetupTestDB setup postgres database with applied app migrations inside docker container
// container bind random port
//
// db *sql.DB   – handle of created DB
// connString string  – connection string to created DB
// close func() – must be called to clean allocated resources. Usually should be deferred immediately
// err – error
func SetupTestDB() (db *sql.DB, connString string, close func() error, err error) {

	var pool *dockertest.Pool

	close = func() error {return nil}
	pool, err = dockertest.NewPool("")
	if err != nil {
		return db, connString, close, err
	}
	pool.MaxWait = time.Second * 30

	connString = os.Getenv("TEST_DB_CONN")
	// If user dont pass TEST_DB_CONN when run test, then we do create isolated db in docker by ourself
	if connString == "" {

		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "12",
			PortBindings: map[dc.Port][]dc.PortBinding{
				"5432": {
					{HostIP: "0.0.0.0", HostPort: "5432"},
				},
			},
			Env: []string{"POSTGRES_PASSWORD=example"},
		})
		if err != nil {
			return db, connString, close, err
		}

		// Sometimes dockertest does not stop container (if test process is killed for example)
		// We configure auto kill container after 5 minutes.
		_  = resource.Expire(5 * 60)

		close = func() error {
			return pool.Purge(resource)
		}
		connString = fmt.Sprintf(
			"postgresql://postgres:example@localhost:%s/postgres?sslmode=disable&search_path=%s",
			resource.GetPort("5432/tcp"), "public",
		)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", connString)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return db, connString, close, err
	}

	migrator, err := migrator_package.NewMigrator(connString)
	if err != nil {
		return db, connString, close, err
	}

	err = migrator.MigrateToLatest()
	if err != nil {
		return db, connString, close, err
	}

	// If user provided its own database for tests then let's
	// down all migrations to cleanup tests run
	if connString != "" {
		close = migrator.DownAllMigrations
	}

	return db, connString, close, err
}
