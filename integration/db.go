package integration

import (
	"database/sql"
	_ "embed"
)

func NewInMemDBConn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "vbs.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
