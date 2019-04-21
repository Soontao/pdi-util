package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

// Version string, in release version
// This variable will be overwrited by complier
var Version = "SNAPSHOT"

// AppName of this application
var AppName = "PDI Util"

// AppUsage of this application
var AppUsage = "A Command Line Tool for SAP PDI"

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Name = AppName
	app.Usage = AppUsage
	app.EnableBashCompletion = true
	app.Flags = globalFlags
	app.Commands = []cli.Command{
		commandSource,
		commandCheck,
		commandSolution,
		commandsPackage,
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
