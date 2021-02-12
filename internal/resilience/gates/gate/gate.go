package gate

// Gate allows the switching of branches for the code to run.
type Gate interface {

	// Switch the branch before running.
	Switch() bool

	// Run the branch.
	Run() error
}
