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
	Description: "Deploy solution without package file, it will download current version package from source tenant, and deploy it to target tenant. Please ensure source & target tenant both are in same release version.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOURCE_SOLUTION_NAME",
			Usage:  "The Source Tenant Solution Name",
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

		version := ""

		sourceSolutionStatus := sourceClient.GetSolutionStatus(sourceSolutionID)

		if sourceSolutionStatus.Status == pdiutil.S_STATUS_ASSEMBLED || sourceSolutionStatus.IsCreatingPatch {
			// if current solution is 'Assembled'
			// or in patch creation
			// Download current version
			version = fmt.Sprintf("%v", sourceSolutionStatus.Version)
		} else {
			// or the latest assembled package
			version = fmt.Sprintf("%v", sourceSolutionStatus.Version-1)
		}

		err, assembledPackage := sourceClient.DownloadSolution(sourceSolutionID, version)

		if err != nil {
			panic(err)
		}

		// content is empty
		if assembledPackage == "" {
			panic(fmt.Errorf("Not found solution %v package with version %v", sourceSolutionID, version))
		} else {
			log.Printf("Assembled packaged downloaded from source tenant, version: %v", version)
		}

		log.Println("Uploading assembled packaged to target system")

		// after deploy, the solution must be existed in target tenant
		err = targetClient.DeploySolution(assembledPackage)

		if err != nil {
			// even successful, server sometimes also response error
			log.Printf("Uploaded to target tenant with error: %v", err)
			log.Println("Even successful, xrep server sometimes also will reset connection")
		} else {
			log.Println("Assembled packaged uploaded to target system")
		}

		// wait seconds
		time.Sleep(checkInterval)

		// use source solution name to find target solution
		targetS := targetClient.GetSolutionByIDOrDescription(sourceSolutionDescription)
		targetSolution := targetS.Name
		targetStatus := targetClient.GetSolutionStatus(targetSolution)

		// if not in uploading and not in uploaded
		if !targetStatus.IsRunningUploading() && !targetStatus.IsUploadingSuccessful() {
			panic("Package uploaded but system not processed, please check the log in system")
		}

		// still in background processing
		if targetStatus.IsRunningUploading() {

			log.Println("Package uploaded, system is processing the uploaded package now")

			// wait uploading finished
			for {
				time.Sleep(checkInterval)
				targetStatus = targetClient.GetSolutionStatus(targetSolution)
				if targetStatus.IsUploadingSuccessful() {
					// not in running mode, break loop
					log.Println("Uploading progress finished")
					break
				} else {
					panic("Upload progress finished, but can not do activation")
				}
			}

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
