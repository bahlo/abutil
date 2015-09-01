package abutil

import (
	"database/sql"
	"errors"
	"fmt"
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

func TestRollbackErr(t *testing.T) {
	mockDBContext(t, func(db *sql.DB) {
		sqlmock.ExpectBegin()
		sqlmock.ExpectRollback()

		tx, err := db.Begin()
		if err != nil {
			t.Error(err)
		}

		alt := errors.New("Some alternative error")
		err = RollbackErr(tx, alt)

		if err != alt {
			t.Errorf("Expected RollbackErr to return %v, but got %v", alt, err)
		}
	})
}

func TestRollbackErrFailing(t *testing.T) {
	mockDBContext(t, func(db *sql.DB) {
		rberr := errors.New("Some rollback error")

		sqlmock.ExpectBegin()
		sqlmock.ExpectRollback().
			WillReturnError(rberr)

		tx, err := db.Begin()
		if err != nil {
			t.Error(err)
		}

		err = RollbackErr(tx, errors.New("This should not be used"))
		if err != rberr {
			t.Errorf("Expected RollbackErr to return %v, but got %v", rberr, err)
		}
	})
}

func rollbackDBContext(fn func(*sql.DB)) {
	db, _ := sqlmock.New()
	fn(db)
	db.Close()
}

func RollbackErrExample() {
	insertSomething := func(db *sql.DB) error {
		tx, _ := db.Begin()

		_, err := tx.Exec("INSERT INTO some_table (some_column) VALUES (?)",
			"foobar")
		if err != nil {
			// We now have a one-liner instead of a check every time an error
			// occurs
			return RollbackErr(tx, err)
		}

		_, err = tx.Exec("DROP DATABASE foobar")
		if err != nil {
			return RollbackErr(tx, err)
		}

		return nil
	}

	rollbackDBContext(func(db *sql.DB) {
		fmt.Println(insertSomething(db))
	})
}
