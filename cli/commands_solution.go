package main

import "github.com/urfave/cli"

var commandSolution = cli.Command{
	Name:  "solution",
	Usage: "solution related operations",
	Subcommands: []cli.Command{
		commandSolutionList,
		commandListSolutionFiles,
		commandSolutionStatus,
		commandSolutionStatusWatch,
		commandSolutionDeploy,
	},
}
