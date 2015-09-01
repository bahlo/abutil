package abutil

import (
	"database/sql"
)

// Rollback does a rollback on the transaction and returns either the error
// from the rollback if there was one or the alternative.
// This is useful if you have multiple statments in a row but don't want to
// call rollback and check for errors every time.
func Rollback(tx *sql.Tx, alt error) error {
	if err := tx.Rollback(); err != nil {
		return err
	}

	return alt
}
