package testing

import (
	"os"
	"strconv"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/pkg/errors"
)

// ConnectionInfo returns a environmental version of the connection info using
// values from the env.
func ConnectionInfo() (database.ConnectionInfo, error) {
	rawPort := getEnvOrElse("POSTGRES_PORT", "5432")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return database.ConnectionInfo{}, errors.WithStack(err)
	}
	return database.ConnectionInfo{
		Host:     getEnvOrElse("POSTGRES_HOST", "localhost"),
		Port:     port,
		User:     getEnvOrElse("POSTGRES_USER", "postgres"),
		Password: getEnvOrElse("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnvOrElse("POSTGRES_DB", "test"),
	}, nil
}

func getEnvOrElse(key string, value string) string {
	envValue, ok := os.LookupEnv(key)
	if ok {
		return envValue
	}
	return value
}
