package integration

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var ddl string

//go:embed dml.sql
var dml string

func NewInMemDBConn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "vbs.db")
	if err != nil {
		return nil, err
	}

	// create tables
	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		return nil, err
	}

	// insert mock data
	if _, err := db.ExecContext(context.Background(), dml); err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
