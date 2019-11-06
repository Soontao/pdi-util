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
	pathes := []string{}

	fileAttrs := c.GetSolutionFileAttrs(solutionName)

	for _, file := range fileAttrs {
		if strings.HasSuffix(file.FilePath, ".absl") || strings.HasSuffix(file.FilePath, ".bo") || strings.HasSuffix(file.FilePath, ".xbo") {
			pathes = append(pathes, file.FilePath)
		}
	}

	sources := c.fetchSources(pathes, concurrent)

	for _, file := range sources {
		check := CopyRightHeaderCheckResult{File: fileAttrs[file.XrepPath], FileContent: *file, HaveHeader: c.checkCopyrightHeader(file.Source)}
		rt = append(rt, check)
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
