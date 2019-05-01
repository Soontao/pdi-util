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

		if sourceSolutionStatus.Status == pdiutil.S_STATUS_ASSEMBLED {
			// if current solution is 'Assembled'
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
		time.Sleep(pdiutil.DefaultPackageCheckInterval)

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
				time.Sleep(pdiutil.DefaultPackageCheckInterval)
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

		// wait seconds
		time.Sleep(pdiutil.DefaultPackageCheckInterval)

		targetStatus = targetClient.GetSolutionStatus(targetSolution)

		if !targetStatus.IsRunningActivation() {
			panic("Activate the solution, but seems system not in progress")
		} else {
			log.Println("Activation running now")
		}

		// wait activation finished
		for {
			time.Sleep(pdiutil.DefaultPackageCheckInterval)
			targetStatus = targetClient.GetSolutionStatus(targetSolution)
			// not in activation now
			if !targetStatus.IsRunningActivation() {
				log.Println("Activation finished")
				break
			}
		}

		if targetStatus.Status == pdiutil.S_STATUS_DEPLOYED {
			// if data update finished so quickly, solution status will be 'Deployed' directly
			log.Println("Deployed")
		} else if targetStatus.IsRunningDataUpdate() {
			// if  data update need time
			log.Println("Data update now")

			// wait data update finished
			for {
				time.Sleep(pdiutil.DefaultPackageCheckInterval)
				targetStatus = targetClient.GetSolutionStatus(targetSolution)

				if !targetStatus.IsRunningDataUpdate() {
					log.Println("Data update finished")
					break
				}
			}

		} else {
			panic("Unknown status, target tenant not deployed but also not in data update, please check the target tenant status")
		}

		log.Println("Finished")

	}),
}
