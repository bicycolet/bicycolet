package statement

// Registrar holds all the statements within a store, that are used to prepare
// them for use.
//
// The store can either register them ahead of time or just in time.
type Registrar interface {

	// Add returns a sql statement if it's found, or will prepare on for use.
	Create(string) (Statement, error)
}
