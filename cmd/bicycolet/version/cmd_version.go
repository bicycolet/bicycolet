package version

import (
	"context"
	"flag"

	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/commands"
	"github.com/spoke-d/clui/flagset"
	"github.com/spoke-d/clui/ui"
	"github.com/spoke-d/task/group"
)

const ClientVersion = "0.0.1-alpha"

type versionCommand struct {
	ui      clui.UI
	flagSet *flagset.FlagSet

	clientVersion string
	template      string
}

func versionCommandFn(version string) func(clui.UI) clui.Command {
	return func(ui clui.UI) clui.Command {
		cmd := &versionCommand{
			ui:      ui,
			flagSet: flagset.New("version", flag.ContinueOnError),

			clientVersion: version,
		}
		cmd.init()
		return cmd
	}
}

func (v *versionCommand) init() {
	v.flagSet.StringVar(&v.template, "template", "{{.ClientVersion}}", "template for the version template")
}

func (v *versionCommand) FlagSet() *flagset.FlagSet {
	return v.flagSet
}

func (v *versionCommand) Usages() []string {
	return make([]string, 0)
}

func (v *versionCommand) Help() string {
	return `
Show the current client version.
`
}

func (v *versionCommand) Synopsis() string {
	return "Show client version."
}

func (v *versionCommand) Init([]string, commands.CommandContext) error {
	return nil
}

func (v *versionCommand) Run(g *group.Group) {
	type version struct {
		ClientVersion string
	}

	template := ui.NewTemplate(versionTemplate, ui.OptionFormat(v.template))
	g.Add(func(ctx context.Context) error {
		return v.ui.Output(template, version{
			ClientVersion: v.clientVersion,
		})
	}, commands.Disguard)
}

const versionTemplate = `
Client version: %s
`
