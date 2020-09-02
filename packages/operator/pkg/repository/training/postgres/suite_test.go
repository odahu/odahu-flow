package postgres_test

import (
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"log"
	"os"
	"testing"
)

var (
	db *sql.DB
)

func Wrapper(m *testing.M) int {
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

	return m.Run()
}

func TestMain(m *testing.M) {

	os.Exit(Wrapper(m))

}
