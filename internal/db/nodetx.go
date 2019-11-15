package db

import (
	"github.com/bicycolet/bicycolet/internal/db/database"
)

// NodeTx models a single interaction with a node-local database.
//
// It wraps low-level db.Tx objects and offers a high-level API to fetch and
// update data.
type NodeTx struct {
	tx    database.Tx // Handle to a transaction in the node-level postgres database.
	query Query
}

// NewNodeTx creates a new transaction node with sane defaults
func NewNodeTx(tx database.Tx) *NodeTx {
	return &NodeTx{
		tx: tx,
	}
}
