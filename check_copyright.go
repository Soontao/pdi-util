package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"baliance.com/gooxml/spreadsheet"
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

// CopyRightHeaderCheckResult struct
type CopyRightHeaderCheckResult struct {
	File       *XrepFile
	HaveHeader bool
}

// CheckSolutionCopyrightHeaderAPI content
func (c *PDIClient) CheckSolutionCopyrightHeaderAPI(solutionName string, concurrent int) []CopyRightHeaderCheckResult {
	rt := []CopyRightHeaderCheckResult{}
	checkList := []string{}
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
		file := <-response
		row := CopyRightHeaderCheckResult{}
		row.File = file
		row.HaveHeader = c.checkCopyrightHeader(file.Source)
		rt = append(rt, row)
	}
	bar.Finish()
	return rt
}

// CheckSolutionCopyrightHeaderToFile content
func (c *PDIClient) CheckSolutionCopyrightHeaderToFile(solutionName string, concurrent int, output string) {

	responses := c.CheckSolutionCopyrightHeaderAPI(solutionName, concurrent)

	tableData := [][]string{}

	for _, r := range responses {

		row := []string{shortenPath2(r.File.XrepPath), strconv.FormatBool(r.HaveHeader), r.File.Attributes["~CREATED_BY"], r.File.Attributes["~LAST_CHANGED_BY"]}

		tableData = append(tableData, row)

	}

	ss := spreadsheet.New()

	addSheetTo(ss, "Copyright Header Check Result", []string{"File", "With Copyright header", "Created By", "Changed By"}, tableData)

	ss.SaveToFile(output)

}

// CheckSolutionCopyrightHeader content
func (c *PDIClient) CheckSolutionCopyrightHeader(solutionName string, concurrent int) {
	lostList := []*XrepFile{}

	checkResult := c.CheckSolutionCopyrightHeaderAPI(solutionName, concurrent)

	for _, response := range checkResult {
		if !response.HaveHeader {
			lostList = append(lostList, response.File)
		}
	}

	if len(lostList) > 0 {
		log.Println("Not found copyright header in: (CreatedBy,ChangedBy:FilePath)")
		for _, file := range lostList {
			c.exitCode = c.exitCode + 1
			log.Printf("%s,%s:%s", file.Attributes["~CREATED_BY"], file.Attributes["~LAST_CHANGED_BY"], shortenPath2(file.XrepPath))
		}
		log.Printf("Totally %d files (of %d) lost copyright header", len(lostList), len(checkResult))
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
		cli.StringFlag{
			Name:   "fileoutput, f",
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name",
		},
	},
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		solutionName := context.String("solution")
		concurrent := context.Int("concurrent")
		output := context.String("fileoutput")
		if output == "" {
			pdiClient.CheckSolutionCopyrightHeader(solutionName, concurrent)
		} else {
			pdiClient.CheckSolutionCopyrightHeaderToFile(solutionName, concurrent, output)
		}

	}),
}
