package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

// Version string, in release version, this variable will be overwrite by complier
var Version = "SNAPSHOT"

// PDIAction wrapper
func PDIAction(action func(pdiClient *PDIClient, c *cli.Context)) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		username := c.GlobalString("username")
		password := c.GlobalString("password")
		hostname := c.GlobalString("hostname")
		hostname = strings.TrimPrefix(hostname, "https://") // remove hostname schema
		pdiClient := NewPDIClient(username, password, hostname)
		action(pdiClient, c) // do process
		if pdiClient.exitCode > 0 {
			// if error happened, change exit code
			return fmt.Errorf("finished %s with error", c.Command.FullName())
		}
		return nil
	}
}

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Name = "PDI Util"
	app.Usage = "A cli util for SAP PDI"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "username, u",
			EnvVar: "PDI_USER",
			Usage:  "The PDI Development User",
		},
		cli.StringFlag{
			Name:   "password, p",
			EnvVar: "PDI_PASSWORD",
			Usage:  "The PDI Development User Password",
		},
		cli.StringFlag{
			Name:   "hostname, t",
			EnvVar: "PDI_TENANT_HOST",
			Usage:  "The PDI Tenant host",
		},
		cli.StringFlag{
			Name:   "output,o",
			EnvVar: "OUPUT",
			Usage:  "output file name",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "source",
			Usage: "source code related operations",
			Subcommands: []cli.Command{
				{
					Name:  "download",
					Usage: "download all files in a solution",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "solution, s",
							EnvVar: "SOLUTION_NAME",
							Usage:  "The PDI Solution Name",
						},
						cli.StringFlag{
							Name:   "output, o",
							EnvVar: "OUTPUT",
							Value:  "output",
							Usage:  "Output directory",
						},
						cli.IntFlag{
							Name:   "concurrent, c",
							EnvVar: "DOWNLOAD_CONCURRENT",
							Value:  35,
							Usage:  "concurrent goroutines number",
						},
					},
					Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
						solutionName := context.String("solution")
						output := context.String("output")
						concurrent := context.Int("concurrent")
						pdiClient.DownloadAllSourceTo(solutionName, output, concurrent)
					}),
				},
			},
		},
		{
			Name:  "session",
			Usage: "session related operations",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list all sessions",
					Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
						solutionName := context.String("solution")
						output := context.String("output")
						concurrent := context.Int("concurrent")
						pdiClient.DownloadAllSourceTo(solutionName, output, concurrent)
					}),
				},
			},
		},
		commandCheck,
		{
			Name:  "solution",
			Usage: "solution related operations",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list all solutions",
					Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
						pdiClient.ListSolutions()
					}),
				},
				{
					Name:  "files",
					Usage: "list all files in a solution",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "solution, s",
							EnvVar: "SOLUTION_NAME",
							Usage:  "The PDI Solution Name",
						},
					},
					Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
						solutionName := context.String("solution")
						pdiClient.ListSolutionAllFiles(solutionName)
					}),
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

}
