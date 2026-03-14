package app

import (
	"database/sql"
	"net/http"
)

type Application struct {
	Server   *http.Server
	Database *sql.DB
}
