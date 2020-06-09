package postgres_test

import (
	"database/sql"
	"fmt"
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type MigratorRouteSuite struct {
	suite.Suite
	g        *GomegaWithT
	pool     *dockertest.Pool
	resource *dockertest.Resource
	migrator *migrator_package.Migrator
}

func (s *MigratorRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

//func (s *MigratorRouteSuite) TearDownTest() {
//	err := s.migrator.DownAllMigrations()
//	if err != nil {
//		s.Suite.FailNow(fmt.Sprintf("Error while teardown test %v", err))
//	}
//}

func (s *MigratorRouteSuite) SetupSuite() {

	//Init postgres repo

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

	// exponential backoff-retry, becapkg/repository/util/postgres/postgres.go:63:47use the application in the container might not be ready to accept connections yet

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

	s.migrator, err = migrator_package.NewMigrator(connString)
	if err != nil {
		panic("Cannot create migrator")
	}

}

func (s *MigratorRouteSuite) TearDownSuite() {
	s.pool.Purge(s.resource)
}

func TestMigratorRouteSuite(t *testing.T) {
	suite.Run(t, new(MigratorRouteSuite))
}

func (s *MigratorRouteSuite) TestMigrateToLatest() {
	err := s.migrator.MigrateToLatest()
	s.g.Expect(err).NotTo(HaveOccurred())
}
