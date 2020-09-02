package deployment_test

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

var (
	db *sql.DB
	kubeClient client.Client
)

const (
	testNamespace = "default"
)

func testMainWrapper(m *testing.M) int {

	// Setup Test DB

	var closeDB func() error
	var err error
	db, _, closeDB, err = testenvs.SetupTestDB()
	defer func() {
		if err := closeDB(); err != nil {
			log.Print("Error during release test DB resources")
		}
	}()
	if err != nil {
		return -1
	}

	var closeKube func() error
	kubeClient, _, closeKube, _, err = testenvs.SetupTestKube(
		filepath.Join("..", "..", "..", "..", "..", "config", "crds"),
		filepath.Join("..", "..", "..", "..", "..", "hack", "tests", "thirdparty_crds"),
	)
	defer func() {
		if err := closeKube(); err != nil {
			log.Print("Error during release test Kube Environment resources")
		}
	}()
	if err != nil {
		return -1
	}

	return m.Run()

}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}