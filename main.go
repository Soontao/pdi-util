package main

import (
	"log"
	"os"
	"sort"

	"github.com/Soontao/pdi-util/client"

	"github.com/urfave/cli"
)

// Version string
var Version string

// PDIAction wrapper
func PDIAction(action func(pdiClient *client.PDIClient, c *cli.Context)) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		pdiClient := client.NewPDIClient(c.GlobalString("username"), c.GlobalString("password"), c.GlobalString("hostname"))
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
							Usage:  "Output directory",
						},
					},
					Action: PDIAction(func(pdiClient *client.PDIClient, context *cli.Context) {
						solutionName := context.String("solution")
						output := context.String("output")
						if output == "" {
							output = "output"
						}
						pdiClient.DownloadAllSourceTo(solutionName, output)
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
					Action: func(c *cli.Context) error {
						return nil
					},
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
					Action: PDIAction(func(pdiClient *client.PDIClient, context *cli.Context) {
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
					Action: PDIAction(func(pdiClient *client.PDIClient, context *cli.Context) {
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