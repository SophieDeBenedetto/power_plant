package database

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
)

// Connect connects to the database with the given credentials
func Connect() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:@localhost/distributed?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	return db
}
