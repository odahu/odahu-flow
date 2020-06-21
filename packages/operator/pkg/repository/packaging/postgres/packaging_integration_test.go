//
//    Copyright 2020 EPAM Systems
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

package postgres_test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	migrator_package "github.com/odahu/odahu-flow/packages/operator/pkg/database/migrators/postgres"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging/postgres"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"log"
	"os"
	"testing"
	"time"
)

const (
	piEntrypoint    = "/usr/bin/test"
	piNewEntrypoint = "/usr/bin/newtest"
	piImage         = "test:image"
	piID            = "pi1"
)

var (
	connString  string
	piArguments = packaging.JsonSchema{
		Properties: []packaging.Property{
			{
				Name: "argument-1",
				Parameters: []packaging.Parameter{
					{
						Name:  "minimum",
						Value: float64(5),
					},
					{
						Name:  "type",
						Value: "number",
					},
				},
			},
		},
		Required: []string{"argument-1"},
	}
	piTargets = []v1alpha1.TargetSchema{
		{
			Name: "target-1",
			ConnectionTypes: []string{
				string(connection.S3Type),
				string(connection.GcsType),
				string(connection.AzureBlobType),
			},
			Required: false,
		},
		{
			Name: "target-2",
			ConnectionTypes: []string{
				string(connection.DockerType),
			},
			Required: true,
		},
	}
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

func TestPackagingRepository(t *testing.T) {

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic("Cannot connect to DB")
	}

	rep := postgres_repo.PackagingIntegrationRepository{DB: db}

	g := NewGomegaWithT(t)

	created := &packaging.PackagingIntegration{
		ID: piID,
		Spec: packaging.PackagingIntegrationSpec{
			Entrypoint:   piEntrypoint,
			DefaultImage: piImage, Privileged: false,
			Schema: packaging.Schema{
				Targets:   piTargets,
				Arguments: piArguments,
			},
		},
	}

	g.Expect(rep.CreatePackagingIntegration(created)).NotTo(HaveOccurred())

	fetched, err := rep.GetPackagingIntegration(piID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := fetched
	updated.Spec.Entrypoint = piNewEntrypoint
	g.Expect(rep.UpdatePackagingIntegration(updated)).NotTo(HaveOccurred())

	fetched, err = rep.GetPackagingIntegration(piID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.Entrypoint).To(Equal(piNewEntrypoint))

	g.Expect(rep.DeletePackagingIntegration(piID)).NotTo(HaveOccurred())
	_, err = rep.GetPackagingIntegration(piID)
	g.Expect(err).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.NotFoundError{Entity: piID}),
	))

}
