package db

import (
	"database/sql"
	"github.com/go-logr/logr"
)

// Commit or Rollback transaction
func FinishTx(tx *sql.Tx, err error, log logr.Logger) {
	if err == nil {
		if err := tx.Commit(); err != nil {
			log.Error(err, "Error while commit transaction")
		}
	} else {
		if err := tx.Rollback(); err != nil {
			log.Error(err, "Error while rollback transaction")
		}
	}
}
