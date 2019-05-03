package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandPackageAssemble = cli.Command{
	Name:  "assemble",
	Usage: "assemble and download assembled package",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
		cli.StringFlag{
			Name:   "output, o",
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "Output file name, if not set, tool will use a default one",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {
		solution := c.GetSolutionIDByString(ctx.String("solution"))

		log.Println("Start current solution status check")
		// check status
		s := c.GetSolutionStatus(solution)

		if s.NeedCreatePatch {
			panic("Solultion need create patch firstly.")
		}
		if !s.CanActivation && !s.CanAssemble {
			panic(fmt.Sprintf("Solution %v can not do activation in: '%v' status", solution, s.StatusText))
		}
		if !s.IsSplitEnabled {
			panic("You need do 'Enabel Assembly Split' manually in PDI.")
		}
		if s.IsRunningJob {
			panic("Another activation/assemble job is running now.")
		}
		if s.IsCreatingPatch {
			panic(fmt.Sprintf("Solution %v is creating patch now.", solution))
		}

		// do checkout check
		// all files must check in, so that you can assemble
		log.Println("Start locks check")
		locks := c.CheckLockedFilesAPI(solution)

		for _, l := range locks {
			log.Printf("%v is locked by %v(%v) at %v", l.FileName, l.EditByUserID, l.EditBy, l.EditOnDate)
		}

		if len(locks) > 0 {
			panic("Solution or solution files are locked by user, please check in them firstly.")
		}

		// do check BAC file
		log.Println("Start BAC check")

		if outOfDate, errs := c.CheckBacOutOfDate(solution); outOfDate {
			for _, e := range errs {
				log.Println(e.Error())
			}
			panic("Please update your BAC file")
		}

		// do check all WCVs have been assigned
		log.Println("Start Check WCV assignment")

		if r := c.FindUnAssignedWCV(solution); r.UnAssignedWCVCount > 0 {
			for _, u := range r.UnAssignedWCVs {
				log.Printf("Un assigned WCV file: %v", u)
			}
			panic("Please make sure all WCV have been assigned")
		}

		// do backend check
		log.Println("Start backend check")
		checkMessages := c.CheckBackendMessageAPI(solution, 30)
		checkErrorCount := 0

		for _, m := range checkMessages {
			if m.IsError() {
				checkErrorCount++
				log.Printf("%v: %v", m.FileName, m.Message)
			}
		}

		if checkErrorCount > 0 {
			panic("Backend code check failed")
		} else {

		}

		// start activation
		log.Println("Start activation")
		if err := c.ActivationSolution(solution); err != nil {
			panic(err)
		}
		// start assemble
		log.Println("Start assemble")
		if err := c.AssembleSolution(solution); err != nil {
			panic(err)
		}

		// start downlaod
		header := c.GetSolutionStatus(solution)
		solutionName := header.SolutionName
		output := ctx.String("output")

		// the version of download
		downloadVersion := ""

		if downloadVersion == "" {
			// default download current version
			if header.Status == pdiutil.S_STATUS_ASSEMBLED {
				downloadVersion = fmt.Sprintf("%v", header.Version)
			} else {
				// or latest version
				downloadVersion = fmt.Sprintf("%v", header.Version-1)
			}
		}

		if output == "" {
			outputID := header.SolutionID
			if header.OriginSolutionID != "" {
				outputID = header.OriginSolutionID
			}
			output = fmt.Sprintf("%v_V%v(%v).zip", outputID, downloadVersion, header.SolutionName)
		}

		log.Printf("Start download %v(%v)", solutionName, downloadVersion)
		err, content := c.DownloadSolution(solution, downloadVersion)

		if err != nil {
			panic(err)
		} else {
			if len(content) != 0 {
				bytes, _ := base64.StdEncoding.DecodeString(content)
				ioutil.WriteFile(output, bytes, 0644)
				log.Printf("File saved to %v", output)
			} else {
				log.Printf("Not found content, please check your version")
			}
		}

		// start create patch solution
		log.Println("Start create patch solution")

		if err := c.CreatePatch(solution); err != nil {
			panic(err)
		}

		log.Println("Finished")

	}),
}
