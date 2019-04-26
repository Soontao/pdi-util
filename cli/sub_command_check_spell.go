package main

import (
	"log"
	"strings"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandCheckSpelling = cli.Command{
	Name:  "spell",
	Usage: "check english spelling",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
		cli.StringFlag{
			Name:   "token",
			EnvVar: "RAPID_API_TOKEN",
			Value:  pdiutil.DefaultRapidAPIToken,
			Usage:  "The Rapid API Token",
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
		results := c.CheckSpellErrorAPI(solution, ctx.String("token"), ctx.Int("concurrent"))

		if len(results) > 0 {
			log.Println("Some words mis-spelling")
			for _, f := range results {
				log.Printf("File %s mis-spelling words: [ %s ]", f.File.XrepPath, strings.Join(f.ErrorSpellingWords, ", "))
			}
		}
		log.Println("Finished")

	}),
}
