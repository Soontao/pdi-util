package main

import "github.com/urfave/cli"

var commandCheck = cli.Command{
	Name:  "check",
	Usage: "static code check",
	Subcommands: []cli.Command{
		commandCheckAll,
		commandCheckCopyright,
		commandCheckBackend,
		commandCheckTranslation,
		commandCheckNameConvention,
		commandCheckUnAssignedWCV,
		commandCheckLocks,
	},
}
