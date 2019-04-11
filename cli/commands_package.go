package main

import (
	"github.com/urfave/cli"
)

var commandsPackage = cli.Command{
	Name:  "package",
	Usage: "package related commands",
	Subcommands: []cli.Command{
		commandPackageDownload,
	},
}
