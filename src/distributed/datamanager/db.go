package datamanager

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:@localhost/distributed?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
}
