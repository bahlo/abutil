package abutil

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func mockDBContext(t *testing.T, fn func(*sql.DB, sqlmock.Sqlmock)) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	fn(db, mock)
}

func TestRollbackErr(t *testing.T) {
	mockDBContext(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectRollback()

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

	mockDBContext(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rberr := errors.New("Some rollback error")

		mock.ExpectBegin()
		mock.ExpectRollback().
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

func exampleRollbackDBContext(fn func(*sql.DB)) {
	db, mock, _ := sqlmock.New()
	mock.ExpectBegin() // At least exopect the begin statement
	fn(db)
	db.Close()
}

func ExampleRollbackErr() {
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

	exampleRollbackDBContext(func(db *sql.DB) {
		fmt.Println(insertSomething(db))
	})
}
