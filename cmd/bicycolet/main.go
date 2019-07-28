package main

import (
	"fmt"
	"os"

	"github.com/bicycolet/bicycolet/pkg/version"
	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/style"
)

const header = `
██████╗ ██╗ ██████╗██╗   ██╗ ██████╗ ██████╗ ██╗     ███████╗████████╗
██╔══██╗██║██╔════╝╚██╗ ██╔╝██╔════╝██╔═══██╗██║     ██╔════╝╚══██╔══╝
██████╔╝██║██║      ╚████╔╝ ██║     ██║   ██║██║     █████╗     ██║
██╔══██╗██║██║       ╚██╔╝  ██║     ██║   ██║██║     ██╔══╝     ██║
██████╔╝██║╚██████╗   ██║   ╚██████╗╚██████╔╝███████╗███████╗   ██║
╚═════╝ ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═════╝ ╚══════╝╚══════╝   ╚═╝
`

const padding = 2

func main() {
	ui := clui.NewColorUI(clui.NewBasicUI(os.Stdin, os.Stdout))
	ui.OutputColor = style.New(style.FgWhite)
	ui.InfoColor = style.New(style.FgGreen)
	ui.WarnColor = style.New(style.FgYellow)
	ui.ErrorColor = style.New(style.FgRed)

	cli := clui.NewCLI("bicycolet", "0.0.1", header, clui.CLIOptions{
		UI: ui,
	})

	cli.AddCommand("version", NewVersionCmd(ui, version.Version))

	exitCode, err := cli.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode.Code())
}

func exit(ui clui.UI, err string) clui.ExitCode {
	ui.Error(err)
	return clui.ExitCode{
		Code: clui.EPerm,
	}
}
