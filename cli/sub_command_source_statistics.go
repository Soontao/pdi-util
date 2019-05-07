package main

import (
	"log"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandSourceStatistics = cli.Command{
	Name:  "statistics",
	Usage: "statistics code scale",
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
		log.Printf("ABSL File Count: %v", result.ABSLFileCount)
		log.Printf("ABSL Code Lines: %v", result.ABSLCodeLines)
		log.Printf("UI/XUI File Count: %v", result.UIComponentCount)
		log.Printf("BO/XBO Count: %v", result.BOCount)
		log.Printf("Web Service Count: %v", result.WebServicesCount)
		log.Printf("Communication Scenario Count: %v", result.CommunicationScenarioCount)

		log.Println("Finished")
	}),
}
