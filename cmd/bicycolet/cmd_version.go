package main

import (
	"flag"
	"runtime"

	"github.com/bicycolet/bicycolet/internal/exec"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/flagset"
)

type versionCmd struct {
	baseCmd
	address       string
	clientVersion string
}

// NewVersionCmd creates a Command with sane defaults
func NewVersionCmd(ui clui.UI, clientVersion string) clui.Command {
	c := &versionCmd{
		baseCmd: baseCmd{
			ui:      ui,
			flagset: flagset.NewFlagSet("version", flag.ExitOnError),
		},
		clientVersion: clientVersion,
	}
	c.init()
	return c
}

func (c *versionCmd) init() {
	c.baseCmd.init()
	c.flagset.StringVar(&c.address, "address", "127.0.0.1:8080", "address of the api server")
}

// Help should return a long-form help text that includes the command-line
// usage. A brief few sentences explaining the function of the command, and
// the complete list of flags the command accepts.
func (c *versionCmd) Help() string {
	return `
Usage:
  version [flags]
Description:
  Show client and server version as JSON, YAML or Tabular.
  If the server is unreachable the version output will label the
  'server_version' as "unreachable".
Example:
  bicycolet version
  bicycolet version --format=json
`
}

// Synopsis should return a one-line, short synopsis of the command.
// This should be short (50 characters of less ideally).
func (c *versionCmd) Synopsis() string {
	return "Show client and server version."
}

// Run should run the actual command with the given CLI instance and
// command-line arguments. It should return the exit status when it is
// finished.
//
// There are a handful of special exit codes that can return documented
// behavioral changes.
func (c *versionCmd) Run() clui.ExitCode {
	// Logging.
	var logger log.Logger
	{
		logLevel := level.AllowInfo()
		if c.debug {
			logLevel = level.AllowAll()
		}
		logger = NewLogCluiFormatter(c.UI())
		logger = log.With(logger,
			"ts", log.DefaultTimestampUTC,
			"uid", uuid.NewRandom().String(),
		)
		logger = level.NewFilter(logger, logLevel)
	}

	client, err := getClient(c.address, logger)
	if err != nil {
		return exit(c.ui, errors.WithStack(err).Error())
	}

	g := exec.NewGroup()
	exec.Block(g)
	{
		g.Add(func() error {
			serverVersion := "unreachable"
			if client != nil {
				result, err := client.Info().Get()
				if err == nil {
					serverVersion = result.Environment.ServerVersion
				}
			}

			version := struct {
				Client  string `json:"client_version" yaml:"client_version" tab:"client"`
				Server  string `json:"server_version" yaml:"server_version" tab:"server"`
				Runtime string `json:"runtime_version" yaml:"runtime_version" tab:"runtime"`
			}{
				Client:  c.clientVersion,
				Server:  serverVersion,
				Runtime: runtime.Version(),
			}

			return c.Output(version)
		}, func(err error) {
			// ignore
		})
	}
	exec.Interrupt(g)
	if err := g.Run(); err != nil {
		return exit(c.ui, err.Error())
	}

	return clui.ExitCode{}
}
