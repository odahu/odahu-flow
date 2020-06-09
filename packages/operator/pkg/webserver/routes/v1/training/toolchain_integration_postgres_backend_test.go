//
//    Copyright 2019 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package training_test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	mt_postgres_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TIPostgresRouteSuite struct {
	TIGenericRouteSuite
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (s *TIPostgresRouteSuite) SetupSuite() {

	// Init postgres repo

	var err error

	s.pool, err = dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	s.pool.MaxWait = time.Second * 30

	// pulls an image, creates a container based on it and runs it
	s.resource, err = s.pool.RunWithOptions(&dockertest.RunOptions{
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
	if err = s.resource.Expire(5 * 60); err != nil {
		panic("Could not set container expire")
	}

	connString := fmt.Sprintf(
		"postgresql://postgres:example@localhost:%s/postgres?sslmode=disable",
		s.resource.GetPort("5432/tcp"),
	)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet

	var db = &sql.DB{}

	if err := s.pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", connString)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		panic("Cannot connect to DB")
	}

	migrator, err := migrator_package.NewMigrator(connString)
	if err != nil {
		panic("Cannot create migrator")
	}
	err = migrator.MigrateToLatest()
	if err != nil {
		panic("Cannot migrate schema")
	}

	s.mtRepository = mt_postgres_repository.ToolchainRepository{DB: db}
}

func (s *TIPostgresRouteSuite) TearDownSuite() {
	s.pool.Purge(s.resource)
}

func TestTIPostgresRouteSuite(t *testing.T) {
	suite.Run(t, new(TIPostgresRouteSuite))
}
