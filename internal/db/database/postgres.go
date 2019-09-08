package database

import (
	_ "github.com/lib/pq"
)

// DriverName to be used for the database.
func DriverName() string {
	return "postgres"
}
