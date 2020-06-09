package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"log"
	"testing"
	"time"
)

func MainTestWrapper(m *testing.M, connString *string) int {

	// For development purposes: to run some tests on persistent database
	if *connString != "" {
		return m.Run()
	}

	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	pool.MaxWait = time.Second * 10

	// pulls an image, creates a container based on it and runs it
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
		panic("Could not start postgres")
	}
	defer pool.Purge(resource)

	*connString = fmt.Sprintf(
		"postgresql://postgres:example@localhost:%s/postgres?sslmode=disable&search_path=%s",
		resource.GetPort("5432/tcp"), "public",
	)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", *connString)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		panic("Cannot connect to DB")
	}

	migrator, err := migrator_package.NewMigrator(*connString)
	if err != nil {
		panic("Cannot create migrator")
	}
	err = migrator.MigrateToLatest()
	if err != nil {
		panic("Cannot migrate schema")
	}

	return m.Run()
}
