package main

import (
	"log"
	"strings"
	"time"

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

var commandSolutionStatusWatch = cli.Command{
	Name:  "watchstatus",
	Usage: "watch the solution status",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {

		solution := c.GetSolutionIDByString(ctx.String("solution"))
		currentStatus := ""
		currentText := ""
		currentRunningAssembleJob := false
		currentRunningCreatePatchJob := false

		for {

			header := c.GetSolutionStatus(solution)

			// print help text after change
			if header.HelpText != currentText {
				currentText = header.HelpText
				for _, l := range strings.Split(currentText, "\\n") {
					log.Printf("System help text: %v", l)
				}
			}

			if header.IsRunningJob != currentRunningAssembleJob {
				currentRunningAssembleJob = header.IsRunningJob
				if currentRunningAssembleJob {
					log.Println("Now solution is running activation/assemble job.")
				} else {
					log.Println("Now solution activation/assemble job finished.")
				}
			}

			if header.IsCreatingPatch != currentRunningCreatePatchJob {
				currentRunningCreatePatchJob = header.IsCreatingPatch
				if currentRunningCreatePatchJob {
					log.Println("Now solution is creating patch solution.")
				} else {
					log.Println("Now solution patch solution created.")
				}

			}

			if header.StatusText != currentStatus {
				currentStatus = header.StatusText
				log.Printf("Now solution %v status is '%v'.", header.SolutionID, header.StatusText)
			}

			// wait default 20 seconds
			time.Sleep(pdiutil.DefaultPackageCheckInterval * time.Second)

		}
	}),
}
