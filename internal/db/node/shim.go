package node

import (
	"database/sql"
	"database/sql/driver"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/pkg/errors"
)

type databaseIO struct {
	// TODO (Simon): Add metrics here
}

func (databaseIO) Register(driverName string, driver driver.Driver) {
	sql.Register(driverName, driver)
}

func (databaseIO) Drivers() []string {
	return sql.Drivers()
}

func (databaseIO) Open(driverName, dataSourceName string) (database.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return database.ShimDB(db, err)
}
