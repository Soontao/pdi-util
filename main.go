package main

import (
	"log"
	"os"
	"sort"

	"github.com/Soontao/pdi-util/client"

	"github.com/urfave/cli"
)

// PDIAction wrapper
func PDIAction(action func(pdiClient *client.PDIClient)) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		pdiClient := client.NewPDIClient(c.String("username"), c.String("password"), c.String("hostname"))
		action(pdiClient)
		return nil
	}
}

func main() {
	app := cli.NewApp()
	app.Version = "v1-alpha"
	app.Name = "PDI Util"
	app.Usage = "A cli util for SAP PDI"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "user, u",
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
			Action: func(c *cli.Context) error {
				return nil
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
					Action: PDIAction(func(pdiClient *client.PDIClient) {
						pdiClient.ListSolutions()
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
