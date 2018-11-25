package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var reg = regexp.MustCompile(`/\*([^*]|[\r\n]|(\*([^/]|[\r\n])))*(Function|Author)([^*]|[\r\n]|(\*([^/]|[\r\n])))*\*/`)

// checkCopyrightHeader
// make sure all absl & bo have copyright header
// with follow format
/*
	Function: make sure all absl & bo have copyright header
	Author: Theo Sun
	Copyright: ?
*/
func (c *PDIClient) checkCopyrightHeader(code []byte) bool {
	return reg.Match(code)
}

// CheckSolutionCopyrightHeader content
func (c *PDIClient) CheckSolutionCopyrightHeader(solutionName string, concurrent int) {
	checkList := []string{}
	lostList := []string{}
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
	asyncResponses := make([]chan bool, fileCount)
	parallexController := make(chan bool, concurrent)
	bar.Start()
	for idx, task := range checkList {
		asyncResponses[idx] = make(chan bool, 1)
		parallexController <- true
		go func(task string, done chan bool) {
			source := c.DownloadFileSource(task)
			done <- c.checkCopyrightHeader(source)
			<-parallexController
			bar.Increment()
		}(task, asyncResponses[idx])
	}
	for idx, response := range asyncResponses {
		haveHeader := <-response // ensure all goroutines finished
		if !haveHeader {
			lostList = append(lostList, checkList[idx])
		}
	}
	bar.Finish()
	for _, file := range lostList {
		log.Printf("Not found copyright header in: %s", strings.TrimPrefix(file, sourceXrepPrefix))
	}
	log.Printf("Totally %d file (of %d) lost copyright header", len(lostList), fileCount)
}
