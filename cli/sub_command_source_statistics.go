package main

import (
	"log"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandSourceStatistics = cli.Command{
	Name:    "statistics",
	Usage:   "statistics code scale",
	Aliases: []string{"scale"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
		cli.IntFlag{
			Name:   "concurrent, c",
			EnvVar: "CHECK_CONCURRENT",
			Value:  35,
			Usage:  "concurrent goroutines number",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {
		solution := c.GetSolutionIDByString(ctx.String("solution"))
		concurrent := ctx.Int("concurrent")
		log.Println("Start statistics")
		result := c.Statistics(solution, concurrent)
		log.Printf("Solution: %v", result.Solution.Name)
		log.Printf("ABSL File Count: %v (%v lines)", result.ABSLFileCount, result.ABSLCodeLines)
		log.Printf("UI/XUI File Count: %v (%v complexity)", result.UIComponentCount, result.UIComplexity)
		log.Printf("BO/XBO Count: %v (%v fields)", result.BOCount, result.BOFieldsCount)
		log.Printf("Communication Scenario Count: %v (%v WebServices)", result.CommunicationScenarioCount, result.WebServicesCount)
		log.Printf("Code List/BCO Count: %v/%v", result.CodeListCount, result.BCOCount)

		log.Println("Finished")
	}),
}
