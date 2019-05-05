package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandPackageDownload = cli.Command{
	Name:  "download",
	Usage: "download assembled package",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
		cli.Int64Flag{
			Name:   "version, v",
			EnvVar: "SOLUTION_VERSION",
			Value:  -1,
			Usage:  "Target download version, if empty, download the latest assembled version",
		},
		cli.StringFlag{
			Name:   "output, o",
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name, default as PDI download package format",
		},
	},
	Action: PDIAction(func(c *pdiutil.PDIClient, ctx *cli.Context) {
		solution := c.GetSolutionIDByString(ctx.String("solution"))
		header := c.GetSolutionStatus(solution)
		solutionName := header.SolutionName
		solutionID := header.SolutionID
		// use original solution id as file part
		if header.OriginSolutionID != "" {
			solutionID = header.OriginSolutionID
		}
		output := ctx.String("output")
		downloadVersion := ctx.Int64("version")

		if downloadVersion < 0 {
			// default download current version
			if header.Status == pdiutil.S_STATUS_ASSEMBLED || header.IsCreatingPatch {
				downloadVersion = header.Version
			} else {
				// or latest version
				downloadVersion = (header.Version - 1)
			}
		}

		if output == "" {
			output = fmt.Sprintf("%v_V%v(%v).zip", solutionID, downloadVersion, solutionName)
		}

		log.Printf("Start download %v(%v)", solutionName, downloadVersion)
		err, content := c.DownloadSolution(solution, fmt.Sprintf("%v", downloadVersion))

		if err != nil {
			log.Panic(err)
		} else {
			if len(content) != 0 {
				bytes, _ := base64.StdEncoding.DecodeString(content)
				ioutil.WriteFile(output, bytes, 0644)
				log.Println("Finished")
			} else {
				log.Printf("Not found content, please check your version")
			}

		}

	}),
}
