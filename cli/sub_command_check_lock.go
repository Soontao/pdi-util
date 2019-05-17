package main

import (
	"log"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandCheckLocks = cli.Command{
	Name:  "lock",
	Usage: "check solution locks",
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
		log.Println("Start check spelling in UI")
		results := c.CheckLockedFilesAPI(solution)

		if len(results) > 0 {
			for _, lock := range results {
				log.Printf("File %s is locked by %s", lock.FilePath, lock.EditByUserID)
			}
			log.Panicln("Some files are locked")
		}

		log.Println("Finished")

	}),
}
