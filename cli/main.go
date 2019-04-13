package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	pdiutil "github.com/Soontao/pdi-util"

	"github.com/urfave/cli"
)

// Version string, in release version
// This variable will be overwrite by complier
var Version = "SNAPSHOT"

// AppName of this application
var AppName = "PDI Util"

// AppUsage of this application
var AppUsage = "A Command Line Tool for SAP Partner Development IDE"

// PDIAction wrapper
func PDIAction(action func(pdiClient *pdiutil.PDIClient, c *cli.Context)) func(c *cli.Context) error {
	return func(c *cli.Context) error {

		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				log.Println("FATAL error happened, so terminated")
				os.Exit(1)
			}
		}()

		// overwrite here
		username := c.GlobalString("username")
		password := c.GlobalString("password")
		hostname := c.GlobalString("hostname")
		release := c.GlobalString("release")

		hostname = strings.TrimPrefix(hostname, "https://") // remove hostname schema
		hostname = strings.TrimPrefix(hostname, "http://")  // remove hostname schema

		pdiClient, err := pdiutil.NewPDIClient(username, password, hostname, release)

		if err != nil {
			return err
		}

		action(pdiClient, c) // do process

		if pdiClient.GetExitCode() > 0 {
			// if error happened, change exit code
			return fmt.Errorf("finished %s with error", c.Command.FullName())
		}
		return nil
	}
}

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
