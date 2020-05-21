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
	Usage: "activate, assemble, download and craete-patch solution",
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
		cli.BoolFlag{
			Name:   "check, c",
			EnvVar: "CHECK_ONLY",
			Usage:  "Do necessary checks for assembly operation only",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {
		solution := c.GetSolutionIDByString(ctx.String("solution"))
		checkOnly := ctx.Bool("check")

		if checkOnly {
			log.Println("Dry run mode")
		}

		log.Printf("Run checks for solution: '%v'", solution)

		// check status
		s := c.GetSolutionStatus(solution)

		if s.NeedCreatePatch {
			panic("Solultion need create patch firstly.")
		}
		if !s.CanActivation && !s.CanAssemble {
			panic(fmt.Sprintf("Solution %v can not do activation in: '%v' status", solution, s.StatusText))
		}
		if !s.IsSplitEnabled {
			panic("You need do 'Enable Assembly Split' manually in PDI.")
		}
		if s.IsRunningJob {
			panic("Another activation/assemble job is running now.")
		}
		if s.IsCreatingPatch {
			panic(fmt.Sprintf("Solution %v is creating patch now.", solution))
		}

		// do checkout check
		// all files must check in, so that you can assemble
		log.Println("Locks check running")
		locks := c.CheckLockedFilesAPI(solution)

		for _, l := range locks {
			log.Println(l.ToString())
		}

		if len(locks) > 0 {
			panic("Solution or solution files are locked by user, please check in them firstly.")
		}
		log.Println("Locks check finished")

		log.Println("In-active file check running")
		inActiveFiles := c.CheckInActiveFilesAPI(solution)

		if len(inActiveFiles) > 0 {

			for _, inActiveFile := range inActiveFiles {
				log.Printf("File '%v', is in-active", inActiveFile.FilePath)
			}

			panic(fmt.Errorf("%v files are in-active, please activate them by PDI firstly", len(inActiveFiles)))

		}
		log.Println("In-active file check finished")

		// do check BAC file
		log.Println("BAC check running")

		if outOfDate, errs := c.CheckBacOutOfDate(solution); outOfDate {
			for _, e := range errs {
				log.Println(e.Error())
			}
			panic("BAC file is out of date, please open your BAC file and activate it in PDI.")
		}

		log.Println("BAC check finished")

		// do check all WCVs have been assigned
		log.Println("WCV assignment check running")

		if r := c.FindUnAssignedWCV(solution); r.UnAssignedWCVCount > 0 {
			for _, u := range r.UnAssignedWCVs {
				log.Printf("Un assigned WCV file: %v", u)
			}
			panic("Please make sure all WCV have been assigned")
		}

		log.Println("WCV assignment check finished")

		// do backend check
		log.Println("Backend check running")
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
		log.Println("Backend check finished")

		// only do check
		if checkOnly {
			log.Println("Check Finished, all things seems fine")
			return
		}

		if s.CanActivation {
			// start activation
			log.Println("Activation running")

			err := c.ActivationSolution(solution)

			// print server log
			for _, slog := range c.GetSolutionLogs(solution, s.Version) {
				log.Printf(">>> Server Log [%v]: %v", slog.Severity, slog.Text)
			}

			if err != nil {
				panic(err)
			}

			log.Println("Activation finished")

		} else {

			log.Println("WARN: The solution is no require to activate, skipped")

		}

		// start assemble
		log.Println("Assemble running")

		err := c.AssembleSolution(solution)

		// print server log
		for _, slog := range c.GetSolutionLogs(solution, s.Version) {
			log.Printf(">>> Server Log [%v]: %v", slog.Severity, slog.Text)
		}

		if err != nil {
			panic(err)
		}

		log.Println("Assemble finished")

		// start download
		header := c.GetSolutionStatus(solution)
		solutionName := header.SolutionName
		output := ctx.String("output")

		// the version of download
		downloadVersion := ""

		// default download current version
		if header.Status == pdiutil.S_STATUS_ASSEMBLED {
			downloadVersion = fmt.Sprintf("%v", header.Version)
		} else {
			// or latest version
			downloadVersion = fmt.Sprintf("%v", header.Version-1)
		}

		if output == "" {
			outputID := header.SolutionID
			if header.OriginSolutionID != "" {
				outputID = header.OriginSolutionID
			}
			output = fmt.Sprintf("%v_V%v(%v).zip", outputID, downloadVersion, header.SolutionName)
		}

		log.Printf("Downloading package %v(%v)", solutionName, downloadVersion)
		err, content := c.DownloadSolution(solution, downloadVersion)

		if err != nil {
			log.Println("Download failed")
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
		log.Println("Creating patch solution")

		if err := c.CreatePatch(solution); err != nil {
			panic(err)
		}

		log.Println("Finished, everything works fine")

	}),
}
