package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

// Version string, in release version, this variable will be overwrite by complier
var Version = "SNAPSHOT"

// PDIAction wrapper
func PDIAction(action func(pdiClient *PDIClient, c *cli.Context)) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		pdiClient := NewPDIClient(c.GlobalString("username"), c.GlobalString("password"), c.GlobalString("hostname"))
		action(pdiClient, c)
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
							Usage:  "concurrent goroutine number when download from remote",
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
		{
			Name:  "check",
			Usage: "static check",
			Subcommands: []cli.Command{
				{
					Name:      "header",
					Usage:     "check copyright header",
					UsageText: "\nmake sure all absl & bo have copyright header with following format:\n\n/*\n\tFunction: make sure all absl & bo have copyright header\n\tAuthor: Theo Sun\n\tCopyright: ?\n*/",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "solution, s",
							EnvVar: "SOLUTION_NAME",
							Usage:  "The PDI Solution Name",
						},
						cli.IntFlag{
							Name:   "concurrent, c",
							EnvVar: "DOWNLOAD_CONCURRENT",
							Value:  35,
							Usage:  "concurrent goroutine number",
						},
					},
					Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
						solutionName := context.String("solution")
						concurrent := context.Int("concurrent")
						pdiClient.CheckSolutionCopyrightHeader(solutionName, concurrent)
					}),
				},
			},
		},
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
