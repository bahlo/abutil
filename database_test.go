package abutil

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func mockDBContext(t *testing.T, fn func(*sql.DB)) {
	db, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	fn(db)
}

func TestRollback(t *testing.T) {
	mockDBContext(t, func(db *sql.DB) {
		sqlmock.ExpectBegin()
		sqlmock.ExpectRollback()

		tx, err := db.Begin()
		if err != nil {
			t.Error(err)
		}

		alt := errors.New("Some alternative error")
		err = Rollback(tx, alt)

		if err != alt {
			t.Errorf("Expected Rollback to return %v, but got %v", alt, err)
		}
	})
}

func TestRollbackFailing(t *testing.T) {
	mockDBContext(t, func(db *sql.DB) {
		rberr := errors.New("Some rollback error")

		sqlmock.ExpectBegin()
		sqlmock.ExpectRollback().
			WillReturnError(rberr)

		tx, err := db.Begin()
		if err != nil {
			t.Error(err)
		}

		err = Rollback(tx, errors.New("This should not be used"))
		if err != rberr {
			t.Errorf("Expected Rollback to return %v, but got %v", rberr, err)
		}
	})
}
