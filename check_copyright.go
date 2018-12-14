package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var reg = regexp.MustCompile(`/\*([^*]|[\r\n]|(\*([^/]|[\r\n])))*(Function|Author)([^*]|[\r\n]|(\*([^/]|[\r\n])))*\*/`)

// checkCopyrightHeader
// make sure all absl & bo have copyright header
// with follow format
/*
	Function: 'Some function description here'
	Author: Theo Sun
	Copyright: ?
*/
func (c *PDIClient) checkCopyrightHeader(code []byte) bool {
	return reg.Match(code)
}

// CheckSolutionCopyrightHeader content
func (c *PDIClient) CheckSolutionCopyrightHeader(solutionName string, concurrent int) {
	checkList := []string{}
	lostList := []*XrepFile{}
	project := c.GetSolutionFileList(solutionName)
	sourceXrepPrefix := ""
	for _, property := range project.PropertyGroup {
		if property.ProjectSourceFolderinXRep != "" {
			sourceXrepPrefix = property.ProjectSourceFolderinXRep
		}
	}
	for _, group := range project.ItemGroup {
		// Common Files
		for _, content := range group.Content {
			if strings.HasSuffix(content.Include, ".absl") || strings.HasSuffix(content.Include, ".bo") || strings.HasSuffix(content.Include, ".xbo") {
				xrepPath := strings.Replace(filepath.Join(sourceXrepPrefix, content.Include), "\\", "/", -1)
				checkList = append(checkList, xrepPath)
			}
		}
	}

	fileCount := len(checkList)
	log.Printf("Will check %d ABSL/BO/XBO Defination\n", fileCount)
	bar := pb.New(fileCount)
	bar.ShowBar = false
	// > request and download
	asyncResponses := make([]chan *XrepFile, fileCount)
	parallexController := make(chan bool, concurrent)
	bar.Start()
	for idx, task := range checkList {
		asyncResponses[idx] = make(chan *XrepFile, 1)
		parallexController <- true
		go func(task string, done chan *XrepFile) {
			source := c.DownloadFileSource(task)
			done <- source
			<-parallexController
			bar.Increment()
		}(task, asyncResponses[idx])
	}
	for _, response := range asyncResponses {
		file := <-response // ensure all goroutines finished
		if !c.checkCopyrightHeader(file.Source) {
			lostList = append(lostList, file)
		}
	}
	bar.Finish()
	if len(lostList) > 0 {
		log.Println("Not found copyright header in: (CreatedBy,ChangedBy:FilePath)")
		for _, file := range lostList {
			c.exitCode = c.exitCode + 1
			log.Printf("%s,%s:%s", file.Attributes["~CREATED_BY"], file.Attributes["~LAST_CHANGED_BY"], shortenPath(strings.TrimPrefix(file.XrepPath, sourceXrepPrefix)))
		}
		log.Printf("Totally %d files (of %d) lost copyright header", len(lostList), fileCount)
	} else {
		log.Println("Congratulation, all source code have copyright header")
	}

}

var commandCheckCopyright = cli.Command{
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
}
