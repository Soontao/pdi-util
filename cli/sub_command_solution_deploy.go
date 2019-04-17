package main

import (
	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandSolutionDeploy = cli.Command{
	Name:  "status",
	Usage: "view the solution status",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOURCE_SOLUTION_NAME",
			Usage:  "The Source Solution Name",
		},
		cli.StringFlag{
			Name:   "source",
			EnvVar: "SOURCE_TENANT",
			Usage:  "The source tenant",
		},
		cli.StringFlag{
			Name:   "sourceuser",
			EnvVar: "SOURCE_TENANT_USER",
			Usage:  "The source tenant",
		},
		cli.StringFlag{
			Name:   "sourcepassword",
			EnvVar: "SOURCE_TENANT_PASSWORD",
			Usage:  "The source tenant",
		},
		cli.StringFlag{
			Name:   "target",
			EnvVar: "TARGET_TENANT",
			Usage:  "The target tenant",
		},
		cli.StringFlag{
			Name:   "targetuser",
			EnvVar: "TARGET_TENANT_USER",
			Usage:  "The target tenant",
		},
		cli.StringFlag{
			Name:   "targetpassword",
			EnvVar: "TARGET_TENANT_PASSWORD",
			Usage:  "The target tenant",
		},
	},
	Action: PDIAction(func(sourceClient *pdiutil.PDIClient, ctx *cli.Context) {
		targetuser := ctx.String("targetuser")
		targetpassword := ctx.String("targetpassword")
		target := ctx.String("target")
		release := ctx.GlobalString("release")
		_, err := pdiutil.NewPDIClient(targetuser, targetpassword, target, release)

		// create target tenant client failed
		if err != nil {
			panic(err)
		}

	}),
}
