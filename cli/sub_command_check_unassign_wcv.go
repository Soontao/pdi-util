package main

import (
	"log"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandCheckUnAssignedWCV = cli.Command{
	Name: "unassignwcv",
	Usage: "c	heck if un-assigned WoC existed.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {
		solution := c.GetSolutionIDByString(ctx.String("solution"))
		log.Println("Start check unassigned WoC files.")
		result := c.FindUnAssignedWCV(solution)
		log.Printf("WC Count: %v", result.WCCount)
		log.Printf("WCV Count: %v", result.WCVCount)
		log.Printf("Assigned WCV Count: %v", result.AssignedWCVCount)
		log.Printf("Un assigned WCV Count: %v", result.UnAssignedWCVCount)

		if result.UnAssignedWCVCount > 0 {
			for _, u := range result.UnAssignedWCVs {
				log.Printf("Un assigned WCV file: %v", u)
			}
			panic("Some WCV not assigned, so failed")
		} else {
			log.Println("All WCVs have been assigned to WC")
		}

	}),
}
