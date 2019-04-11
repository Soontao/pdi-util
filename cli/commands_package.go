package main

import (
	"github.com/urfave/cli"
)

var commandsPackage = cli.Command{
	Name:  "package",
	Usage: "package related operations",
	Subcommands: []cli.Command{
		commandPackageDownload,
		commandPackageAssemble,
	},
}
