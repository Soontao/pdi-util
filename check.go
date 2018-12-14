package main

import (
	"github.com/urfave/cli"
)

// empty index
// content splited

var commandCheck = cli.Command{
	Name:  "check",
	Usage: "static check",
	Subcommands: []cli.Command{
		commandCheckCopyright,
		commandCheckBackend,
		commandCheckTranslation,
		commandCheckNameConvention,
	},
}
