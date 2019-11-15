package db

import (
	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/bicycolet/bicycolet/internal/db/schema"
)

// NodeTransactioner represents a way to run transaction on the node
type NodeTransactioner interface {

	// Transaction creates a new NodeTx object and transactionally executes the
	// node-level database interactions invoked by the given function. If the
	// function returns no error, all database changes are committed to the
	// node-level database, otherwise they are rolled back.
	Transaction(f func(*NodeTx) error) error
}

// QueryNode represents a local node in a cluster
type QueryNode interface {

	// Open the node-local database object.
	Open(string) error

	// EnsureSchema applies all relevant schema updates to the node-local
	// database.
	//
	// Return the initial schema version found before starting the update, along
	// with any error occurred.
	EnsureSchema(hookFn schema.Hook) (int, error)

	// DB return the current database source.
	DB() database.DB
}

type nodeTxBuilder func(database.Tx) *NodeTx

// Node mediates access to the data stored in the node-local postgres database.
type Node struct {
	transaction Transaction // Handle the transactions to the database
	node        QueryNode
	dir         string // Reference to the directory where the database file lives.
	builder     nodeTxBuilder
}

// Transaction creates a new NodeTx object and transactionally executes the
// node-level database interactions invoked by the given function. If the
// function returns no error, all database changes are committed to the
// node-level database, otherwise they are rolled back.
func (n *Node) Transaction(f func(*NodeTx) error) error {
	return n.transaction.Transaction(n.node.DB(), func(tx database.Tx) error {
		return f(n.builder(tx))
	})
}

// Close the database facade.
func (n *Node) Close() error {
	return n.node.DB().Close()
}
