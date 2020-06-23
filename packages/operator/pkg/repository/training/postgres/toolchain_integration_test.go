package postgres_test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"log"
	"os"
	"testing"
	"time"
)

var (
	connString string
)

const (
	tiID            = "foo"
	tiEntrypoint    = "test-entrypoint"
	tiNewEntrypoint = "new-test-entrypoint"
)

func Wrapper(m *testing.M) int {
	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	pool.MaxWait = time.Second * 30

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

	if err = resource.Expire(5 * 60); err != nil {
		panic("Could not set container expire")
	}

	connString = fmt.Sprintf(
		"postgresql://postgres:example@localhost:%s/postgres?sslmode=disable&search_path=%s",
		resource.GetPort("5432/tcp"), "public",
	)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", connString)
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

	return m.Run()
}

func TestMain(m *testing.M) {

	os.Exit(Wrapper(m))

}

func TestToolchainRepository(t *testing.T) {

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic("Cannot connect to DB")
	}

	tRepo := postgres_repo.ToolchainRepository{DB: db}

	g := NewGomegaWithT(t)

	created := &training.ToolchainIntegration{
		ID: tiID,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			Entrypoint: tiEntrypoint,
		},
	}

	g.Expect(tRepo.CreateToolchainIntegration(created)).NotTo(HaveOccurred())

	g.Expect(tRepo.CreateToolchainIntegration(created)).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.AlreadyExistError{Entity: tiID}),
	))

	fetched, err := tRepo.GetToolchainIntegration(tiID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := &training.ToolchainIntegration{
		ID: tiID,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			Entrypoint: tiNewEntrypoint,
		},
	}
	g.Expect(tRepo.UpdateToolchainIntegration(updated)).NotTo(HaveOccurred())

	fetched, err = tRepo.GetToolchainIntegration(tiID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.Entrypoint).To(Equal(tiNewEntrypoint))

	tis, err := tRepo.GetToolchainIntegrationList()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(tis)).To(Equal(1))

	g.Expect(tRepo.DeleteToolchainIntegration(tiID)).NotTo(HaveOccurred())
	_, err = tRepo.GetToolchainIntegration(tiID)
	g.Expect(err).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.NotFoundError{Entity: tiID}),
	))

}
