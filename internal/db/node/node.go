package node

import (
	"database/sql/driver"
	"path/filepath"
	"sync"
	"time"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/bicycolet/bicycolet/internal/db/schema"
	"github.com/bicycolet/bicycolet/internal/fsys"
	"github.com/pkg/errors"
)

// DatabaseRegistrar represents a way to register and un-register drivers with
// associated name.
type DatabaseRegistrar interface {
	// Register makes a database driver available by the provided name.
	// If Register is called twice with the same name or if driver is nil,
	// it panics.
	Register(name string, driver driver.Driver)

	// Drivers returns a sorted list of the names of the registered drivers.
	Drivers() []string
}

// DatabaseOpener represents a way to open a database source
type DatabaseOpener interface {
	// Open opens a database specified by its database driver name and a
	// driver-specific data source name, usually consisting of at least a
	// database name and connection information.
	Open(driverName, dataSourceName string) (database.DB, error)
}

// DatabaseIO captures all the input/output (IO) of the database in one logical
// interface.
type DatabaseIO interface {
	DatabaseRegistrar
	DatabaseOpener
}

// SchemaProvider defines a factory for creating schemas
type SchemaProvider interface {

	// Schema creates a new Schema that captures the schema of a database in
	// terms of a series of ordered updates.
	Schema() Schema

	// Updates returns the schema updates that is required for the updating of
	// the database.
	Updates() []schema.Update
}

// Schema captures the schema of a database in terms of a series of ordered
// updates.
type Schema interface {

	// File extra queries from a file. If the file is exists, all SQL queries in it
	// will be executed transactionally at the very start of Ensure(), before
	// anything else is done.
	File(string)
	// Hook instructs the schema to invoke the given function whenever a update is
	// about to be applied. The function gets passed the update version number and
	// the running transaction, and if it returns an error it will cause the schema
	// transaction to be rolled back. Any previously installed hook will be
	// replaced.
	Hook(schema.Hook)

	// Ensure makes sure that the actual schema in the given database matches the
	// one defined by our updates.
	//
	// All updates are applied transactionally. In case any error occurs the
	// transaction will be rolled back and the database will remain unchanged.
	Ensure(database.DB) (int, error)
}

// Node represents a local node in a cluster
type Node struct {
	database       database.DB
	databasePath   string
	databaseIO     DatabaseIO
	schemaProvider SchemaProvider
	fileSystem     fsys.FileSystem
	openTimeout    time.Duration
	once           sync.Once
}

// New creates a cluster ensuring that sane defaults are employed.
func New(fileSystem fsys.FileSystem) *Node {
	return &Node{
		databaseIO: databaseIO{},
		schemaProvider: &schemaProvider{
			fileSystem: fileSystem,
		},
		fileSystem: fileSystem,
	}
}

// Open the node-local database object.
func (n *Node) Open(path string, connectionInfo database.ConnectionInfo) error {
	// These are used to tune the transaction BEGIN behavior instead of using the
	// similar "locking_mode" pragma (locking for the whole database connection).
	db, err := n.databaseIO.Open(database.DriverName(), connectionInfo.String())
	if err != nil {
		return errors.WithStack(err)
	}
	if err := db.Ping(); err != nil {
		return errors.WithStack(err)
	}

	n.database = db
	n.databasePath = path

	return nil
}

// EnsureSchema applies all relevant schema updates to the node-local
// database.
//
// Return the initial schema version found before starting the update, along
// with any error occurred.
func (n *Node) EnsureSchema(hookFn schema.Hook) (int, error) {
	ctx := &context{}

	schema := n.schemaProvider.Schema()
	schema.File(filepath.Join(n.databasePath, "patch.local.sql"))
	schema.Hook(func(version int, tx database.Tx) error {
		err := hook(ctx, n.fileSystem, hookFn, n.databasePath, version, tx)
		return errors.WithStack(err)
	})
	return schema.Ensure(n.database)
}

// DB return the current database source.
func (n *Node) DB() database.DB {
	return n.database
}

type context struct {
	backupDone bool
}

func hook(ctx *context, fsys fsys.FileSystem, hook schema.Hook, dir string, version int, tx database.Tx) error {
	if !ctx.backupDone {
		ctx.backupDone = true
	}

	// Run the given hook only against actual update versions, not
	// when a custom query file is passed (signaled by version == -1).
	if hook != nil && version != -1 {
		if err := hook(version, tx); err != nil {
			return err
		}
	}
	return nil
}
