package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var reg = regexp.MustCompile(`/\*([^*]|[\r\n]|(\*([^/]|[\r\n])))*(Function|Author)([^*]|[\r\n]|(\*([^/]|[\r\n])))*\*/`)

// maybe these rules can be loaded by external
var rules = map[string]string{
	// extension: prefix
	"xbo":              "BOE",
	"bo":               "BO",
	"csd":              "CS",
	"bco":              "BCO",
	"codelist":         "CLDT",
	"EC.uicomponent":   "EC",
	"WCF.uiwoc":        "UI",
	"WCVIEW.uiwocview": "UI",
	"library":          "RL",
}

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

func ensureFileNameConvention(filePath string) (bool, string) {

	filePathSlices := strings.Split(filePath, "\\")
	fileName := filePathSlices[len(filePathSlices)-1]

	legal := true
	correctFileName := fileName
	fileExtensionArray := strings.SplitN(fileName, ".", 2)
	fileNameWithoutExtension := fileExtensionArray[0]
	fileExtension := fileExtensionArray[1]
	fileNameSlice := strings.SplitN(fileNameWithoutExtension, "_", 2)
	fileNamePrefix := ""
	fileNameWithoutPrefix := fileNameWithoutExtension

	if len(fileNameSlice) == 2 {
		fileNamePrefix = fileNameSlice[0]
		fileNameWithoutPrefix = fileNameSlice[1]
	}

	correctPrefix := rules[fileExtension]

	if correctPrefix != "" && fileNamePrefix != correctPrefix {
		legal = false
		correctFileName = correctPrefix + "_" + fileNameWithoutPrefix + "." + fileExtension
	}

	return legal, correctFileName
}

// CheckNameConvention of the solution
func (c *PDIClient) CheckNameConvention(solution string) {
	count := 0
	project := c.GetSolutionFileList(solution)
	for _, group := range project.ItemGroup {
		for _, bcset := range group.BCSet {
			includePath := strings.SplitN(bcset.Include, "\\", 3)[2]
			correct, correcetName := ensureFileNameConvention(includePath)
			if !correct {
				count = count + 1
				log.Printf("The name of file %s should be %s\n", includePath, correcetName)
			}
		}
		for _, content := range group.Content {
			includePath := content.Include
			correct, correcetName := ensureFileNameConvention(includePath)
			if !correct {
				count = count + 1
				log.Printf("The name of file %s should be %s\n", includePath, correcetName)
			}

		}
	}
	log.Printf("finished, name convension error count: %d", count)
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
			log.Printf("%s,%s:%s", file.Attributes["~CREATED_BY"], file.Attributes["~LAST_CHANGED_BY"], strings.TrimPrefix(file.XrepPath, sourceXrepPrefix))
		}
		log.Printf("Totally %d files (of %d) lost copyright header", len(lostList), fileCount)
	} else {
		log.Println("Congratulation, all source code have copyright header")
	}

}
