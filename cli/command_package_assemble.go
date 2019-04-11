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
			Name:   "version, v",
			EnvVar: "SOLUTION_VERSION",
			Usage:  "Target download version",
		},
		cli.StringFlag{
			Name:   "output, o",
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {
		solution := c.GetSolutionIDByString(ctx.String("solution"))

		// check status
		s := c.GetSolutionStatus(solution)
		if s.NeedCreatePatch {
			log.Panicln("Solultion need create patch firstly.")
			return
		}
		if !s.CanActivation && !s.CanAssemble {
			log.Panicf("Solution %v can not do activation in: %v status", solution, s.StatusText)
			return
		}

		// start activation
		log.Println("Start activation")
		if err := c.ActivationSolution(solution); err != nil {
			log.Println(err)
			return
		}
		// start assemble
		log.Println("Start assemble")
		if err := c.AssembleSolution(solution); err != nil {
			log.Println(err)
			return
		}

		// start downlaod
		header := c.GetSolutionStatus(solution)
		solutionName := header.SolutionName
		output := ctx.String("output")
		downloadVersion := ctx.String("version")

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
			output = fmt.Sprintf("%v(V%v).zip", solutionName, downloadVersion)
		}

		log.Printf("Start download %v(%v)", solutionName, downloadVersion)
		err, content := c.DownloadSolution(solution, downloadVersion)

		if err != nil {
			log.Panic(err)
		} else {
			if len(content) != 0 {
				bytes, _ := base64.StdEncoding.DecodeString(content)
				ioutil.WriteFile(output, bytes, 0644)
				log.Println("Finished")
			} else {
				log.Printf("Not found content.")
			}

		}
		// start create patch solution
		log.Println("Start create patch solution")
		if err := c.CreatePatch(solution); err != nil {
			log.Println(err)
			return
		}

		log.Println("Finished")

	}),
}
