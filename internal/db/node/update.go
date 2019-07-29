package node

import (
	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/bicycolet/bicycolet/internal/db/schema"
	"github.com/bicycolet/bicycolet/internal/fsys"
)

type schemaProvider struct {
	fileSystem fsys.FileSystem
}

func (s schemaProvider) Schema() Schema {
	schema := schema.New(s.fileSystem, s.Updates())
	schema.Fresh(freshSchema)
	return schema
}

func (s schemaProvider) Updates() []schema.Update {
	return []schema.Update{
		updateFromV0,
	}
}

func updateFromV0(tx database.Tx) error {
	stmt := `
CREATE TABLE config (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	key VARCHAR(255) NOT NULL,
	value TEXT,
	UNIQUE (key)
);
CREATE TABLE patches (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name VARCHAR(255) NOT NULL,
	applied_at DATETIME NOT NULL,
	UNIQUE (name)
);
`
	_, err := tx.Exec(stmt)
	return err
}
