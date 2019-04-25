package main

import (
	"github.com/urfave/cli"
)

var commandCheckSpelling = cli.Command{
	Name:  "spell",
	Usage: "check english spelling",
}
