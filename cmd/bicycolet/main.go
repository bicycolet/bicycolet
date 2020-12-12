package main

import (
	"fmt"
	"os"

	"github.com/bicycolet/bicycolet/cmd/bicycolet/version"
	"github.com/spoke-d/clui"
	"github.com/spoke-d/clui/autocomplete/fsys"
)

func main() {
	fsys := fsys.NewLocalFileSystem()

	cli := clui.New("bicycolet", version.ClientVersion, header, clui.OptionFileSystem(fsys))

	version.Register(cli)

	code, err := cli.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code.Code())
}

const header = `
██████╗ ██╗ ██████╗██╗   ██╗ ██████╗ ██████╗ ██╗     ███████╗████████╗
██╔══██╗██║██╔════╝╚██╗ ██╔╝██╔════╝██╔═══██╗██║     ██╔════╝╚══██╔══╝
██████╔╝██║██║      ╚████╔╝ ██║     ██║   ██║██║     █████╗     ██║
██╔══██╗██║██║       ╚██╔╝  ██║     ██║   ██║██║     ██╔══╝     ██║
██████╔╝██║╚██████╗   ██║   ╚██████╗╚██████╔╝███████╗███████╗   ██║
╚═════╝ ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═════╝ ╚══════╝╚══════╝   ╚═╝
`
