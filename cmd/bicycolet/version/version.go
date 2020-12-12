package version

import (
	"github.com/spoke-d/clui"
)

// Registry defines a use in site interface for registering commands.
type Registry interface {
	// Add a command.
	Add(string, clui.CommandFn) error
}

// Register adds commands to the registry.
func Register(cli Registry) {
	cli.Add("version", versionCommandFn(ClientVersion))
}
