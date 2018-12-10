package main

import (
	"github.com/urfave/cli"
)

// empty index
// content splited

var commandCheck = cli.Command{
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
					EnvVar: "CHECK_CONCURRENT",
					Value:  35,
					Usage:  "concurrent goroutines number",
				},
			},
			Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
				solutionName := context.String("solution")
				concurrent := context.Int("concurrent")
				pdiClient.CheckSolutionCopyrightHeader(solutionName, concurrent)
			}),
		},
		commandCheckBackend,
		{
			Name:  "translation",
			Usage: "do translation check",
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
				cli.StringFlag{
					Name:   "language, l",
					EnvVar: "LANGUAGE",
					Value:  "Chinese",
					Usage:  "target language to check",
				},
			},
			Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
				solutionName := context.String("solution")
				concurrent := context.Int("concurrent")
				language := context.String("language")
				pdiClient.CheckTranslation(solutionName, concurrent, language)
			}),
		},
		{
			Name:      "name",
			Usage:     "check name convension",
			UsageText: "check the name convension of source code",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "solution, s",
					EnvVar: "SOLUTION_NAME",
					Usage:  "The PDI Solution Name",
				},
			},
			Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
				solutionName := context.String("solution")
				pdiClient.CheckNameConvention(solutionName)
			}),
		},
	},
}
