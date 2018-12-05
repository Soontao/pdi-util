package main

import (
	"log"
	"strings"
)

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

func shortenPath(path string) string {
	rt := path
	pathArray := strings.Split(path, "\\")
	pathArrayLen := len(pathArray)

	if pathArrayLen > 1 {
		for idx := 0; idx < pathArrayLen-1; idx++ {
			if pathArray[idx] != ".." {
				// pick first char
				pathArray[idx] = pathArray[idx][0:1]
			}
		}
		rt = strings.Join(pathArray, "\\")
	}

	return rt
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
				log.Printf("The name should be %s of file %s\n", correcetName, shortenPath(includePath))
			}
		}
		for _, content := range group.Content {
			includePath := content.Include
			correct, correcetName := ensureFileNameConvention(includePath)
			if !correct {
				count = count + 1
				log.Printf("The name should be %s of file %s\n", correcetName, shortenPath(includePath))
			}

		}
	}
	c.exitCode = count
	log.Printf("name convension error count: %d", count)
}
