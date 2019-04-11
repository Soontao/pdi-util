package main

import (
	"log"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandSolutionStatus = cli.Command{
	Name:  "status",
	Usage: "view the solution status",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {
		solution := c.GetSolutionIDByString(ctx.String("solution"))
		header := c.GetSolutionStatus(solution)
		log.Printf("Solution ID:\t%v", header.SolutionID)
		log.Printf("Solution Name:\t%v", header.SolutionName)
		log.Printf("Solution Status:\t%v", header.StatusText)

		log.Printf("Solution Enabled:\t%v", header.Enabled)
		log.Printf("Solution Version:\t%v", header.Version)
		log.Printf("Solution UpdatedAt:\t%v", header.ChangeDateTime)
	}),
}
