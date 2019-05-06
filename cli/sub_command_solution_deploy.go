package main

import (
	"fmt"
	"log"
	"time"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandSolutionDeploy = cli.Command{
	Name:        "deploy",
	Usage:       "deploy solution",
	Description: "Deploy solution without package file, it will download assembled package from source tenant, and deploy it to target tenant. Please ensure source & target tenant both are in same release version.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOURCE_SOLUTION_NAME",
			Usage:  "The source tenant solution name, recommend to use the readable solution description",
		},
		cli.Int64Flag{
			Name:   "version",
			EnvVar: "SOURCE_SOLUTION_VERSION",
			Value:  -1,
			Usage:  "The specific solution version to deploy, if empty, use the latest assembled version",
		},
		cli.StringFlag{
			Name:   "target",
			EnvVar: "TARGET_TENANT",
			Usage:  "The target tenant",
		},
		cli.StringFlag{
			Name:   "targetuser",
			EnvVar: "TARGET_TENANT_USER",
			Usage:  "The target tenant user",
		},
		cli.StringFlag{
			Name:   "targetpassword",
			EnvVar: "TARGET_TENANT_PASSWORD",
			Usage:  "The target tenant user password",
		},
	},
	Action: PDIAction(func(sourceClient *pdiutil.PDIClient, ctx *cli.Context) {

		// check solution status every 20 seconds
		checkInterval := time.Second * 20

		targetClient, err := pdiutil.NewPDIClient(
			ctx.String("targetuser"),
			ctx.String("targetpassword"),
			ctx.String("target"),
			ctx.GlobalString("release"),
		)

		// create target tenant client failed
		if err != nil {
			panic(err)
		}

		s := sourceClient.GetSolutionByIDOrDescription(ctx.String("solution"))
		sourceSolutionID := s.Name
		sourceSolutionDescription := s.Description

		sourceSolutionStatus := sourceClient.GetSolutionStatus(sourceSolutionID)

		version := ctx.Int64("version")

		// if not specific external version
		if version < 0 {
			version = sourceSolutionStatus.GetSolutionLatestAssembledVersion()
		}

		log.Printf("Downloading the assembled package from source tenant, version: %v", version)

		err, assembledPackage := sourceClient.DownloadSolution(sourceSolutionID, fmt.Sprintf("%v", version))

		if err != nil {
			panic(err)
		}

		// content is empty
		if assembledPackage == "" {
			panic(fmt.Errorf("Not found solution %v package with version %v", sourceSolutionID, version))
		}

		log.Println("Uploading assembled package to target system")

		// after deploy, the solution must be existed in target tenant
		if err = targetClient.DeploySolution(assembledPackage); err != nil {
			panic(err)
		} else {
			log.Println("Uploaded")
		}

		// wait seconds
		time.Sleep(checkInterval)

		// use source solution name to find target solution
		targetS := targetClient.GetSolutionByIDOrDescription(sourceSolutionDescription)
		targetSolution := targetS.Name
		targetStatus := targetClient.GetSolutionStatus(targetSolution)

		// if not in uploading and not in uploaded
		if !targetStatus.IsUploadingSuccessful() {
			panic("Package uploaded but system not processed, please check the log in system")
		}

		// activate target solution
		if err = targetClient.ActivateDeployedSolution(targetSolution); err != nil {
			panic(err)
		}

		log.Println("Activation triggered, wait solution activation & data update now")

		for {

			time.Sleep(checkInterval)

			targetStatus = targetClient.GetSolutionStatus(targetSolution)

			if targetStatus.Status == pdiutil.S_STATUS_DEPLOYED {
				log.Println("Deployed")
				break
			}

		}

		log.Println("Finished")

	}),
}
