package pdiutil

import (
	"log"
	"regexp"
	"strings"
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
	File        XrepFileAttrs
	FileContent XrepFile
	HaveHeader  bool
}

// CheckSolutionCopyrightHeaderAPI content
func (c *PDIClient) CheckSolutionCopyrightHeaderAPI(solutionName string, concurrent int) []CopyRightHeaderCheckResult {
	rt := []CopyRightHeaderCheckResult{}

	files := c.GetSolutionFileAttrs(solutionName)

	for _, file := range files {
		if strings.HasSuffix(file.FilePath, ".absl") || strings.HasSuffix(file.FilePath, ".bo") || strings.HasSuffix(file.FilePath, ".xbo") {
			rt = append(rt, CopyRightHeaderCheckResult{File: file})
		}
	}

	fileCount := len(rt)
	// > request and download
	asyncResponses := make([]chan bool, fileCount)
	parallexController := make(chan bool, concurrent)
	for idx, task := range rt {
		asyncResponses[idx] = make(chan bool, 1)
		parallexController <- true
		go func(task CopyRightHeaderCheckResult, done chan bool) {
			task.FileContent = *(c.DownloadFileSource(task.File.FilePath))
			done <- true
			<-parallexController
		}(task, asyncResponses[idx])
	}
	// await all go-routines finished
	for _, response := range asyncResponses {
		<-response
	}

	for _, check := range rt {
		fileContent := check.FileContent
		check.HaveHeader = c.checkCopyrightHeader(fileContent.Source)
	}
	return rt
}

// CheckSolutionCopyrightHeader content
func (c *PDIClient) CheckSolutionCopyrightHeader(solutionName string, concurrent int) {
	lostList := []XrepFile{}

	checkResult := c.CheckSolutionCopyrightHeaderAPI(solutionName, concurrent)

	for _, response := range checkResult {
		if !response.HaveHeader {
			lostList = append(lostList, response.FileContent)
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
