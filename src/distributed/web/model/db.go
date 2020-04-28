package model

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/sophiedebenedetto/power_plant/src/distributed/database"
)

var db *sql.DB

func init() {
	db = database.Connect()
}
