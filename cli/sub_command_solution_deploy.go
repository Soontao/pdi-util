package main

import (
	"fmt"
	"log"
	"time"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandSolutionDeploy = cli.Command{
	Name:  "deploy",
	Usage: "deploy solution to target tenant",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOURCE_SOLUTION_NAME",
			Usage:  "The Source Solution Name",
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

		// if current solution is 'Assembled'
		// Download current version
		if sourceSolutionStatus.Status == pdiutil.S_STATUS_ASSEMBLED {
			version = fmt.Sprintf("%v", sourceSolutionStatus.Version)
		} else {
			// or latest version
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
			log.Printf("Assmebled packaged downloaded from source tenant")
		}

		// after deploy, the solution must be existed in target tenant
		err = targetClient.DeploySolution(assembledPackage)

		log.Println("Assmebled packaged uploaded to target system")

		if err != nil {
			// even successful, server sometimes also response error
			log.Println(err)
		}

		// use source solution name to find target solution
		targetS := targetClient.GetSolutionByIDOrDescription(sourceSolutionDescription)
		targetSolution := targetS.Name
		targetStatus := targetClient.GetSolutionStatus(targetSolution)

		if !targetStatus.IsRunningUploading() {
			panic("Package uploaded but system not processed, please check log in system")
		}

		log.Println("Package uploaded, system is processing the uploaded package")

		// wait uploading finished
		for {

			time.Sleep(pdiutil.DefaultPackageCheckInterval)
			targetStatus = targetClient.GetSolutionStatus(targetSolution)
			if !targetStatus.IsRunningUploading() {
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
		err = targetClient.ActivateDeployedSolution(targetSolution)

		if err != nil {
			panic(err)
		}

		targetStatus = targetClient.GetSolutionStatus(targetSolution)

		if !targetStatus.IsRunningActivation() {
			panic("Activate the solution, but seems system not in progress")
		}

		log.Println("Activate the solution now")
		// wait activation finished
		for {
			time.Sleep(pdiutil.DefaultPackageCheckInterval)
			targetStatus = targetClient.GetSolutionStatus(targetSolution)
			if !targetStatus.IsRunningActivation() {
				log.Println("Activation finished")
				break
			}
		}

		log.Println("Data update now")
		// wait data update finished
		for {
			time.Sleep(pdiutil.DefaultPackageCheckInterval)
			targetStatus = targetClient.GetSolutionStatus(targetSolution)
			if !targetStatus.IsRunningDateUpdate() {
				log.Println("Data update finished")
				break
			}
		}

		log.Println("Finished")

	}),
}
