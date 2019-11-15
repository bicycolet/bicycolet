package testing

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/pkg/errors"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const dbNameLength = 40

// RandomTableName generates a random DB name
func RandomTableName() string {
	b := make([]byte, dbNameLength)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := dbNameLength-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

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
