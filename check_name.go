package pdiutil

import (
	"log"
	"strconv"
	"strings"

	"baliance.com/gooxml/spreadsheet"
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

func shortenPath2(path string) string {
	rt := path
	path = strings.Replace(path, "\\", "/", -1)
	pathArray := strings.Split(path, "/")
	pathArrayLen := len(pathArray)

	if pathArrayLen > 1 {
		for idx := 0; idx < pathArrayLen-1; idx++ {
			if pathArray[idx] != ".." && pathArray[idx] != "" {
				// pick first char
				pathArray[idx] = pathArray[idx][0:1]
			}
		}
		rt = strings.Join(pathArray, "/")
	}

	return rt
}

// ShortenPath2 exported
var ShortenPath2 = shortenPath2

func shortenPath(path string) string {
	rt := path
	pathArray := strings.Split(path, "\\")
	pathArrayLen := len(pathArray)

	if pathArrayLen > 1 {
		for idx := 0; idx < pathArrayLen-1; idx++ {
			if pathArray[idx] != ".." && pathArray[idx] != "" {
				// pick first char
				pathArray[idx] = pathArray[idx][0:1]
			}
		}
		rt = strings.Join(pathArray, "\\")
	}

	return rt
}

func ensureFileNameConvention(filePath string) (bool, string) {

	filePathSlices := strings.Split(filePath, "/")
	fileName := filePathSlices[len(filePathSlices)-1]

	legal := true
	correctFileName := fileName
	fileExtensionArray := strings.SplitN(fileName, ".", 2)
	// the file path without extension
	if len(fileExtensionArray) < 2 {
		return true, filePath
	}
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

// NameConventionCheckResult type
type NameConventionCheckResult struct {
	Correct     bool
	IncludePath string
	CorrectName string
	File        XrepFileAttrs
}

// CheckNameConventionAPI of the solution
func (c *PDIClient) CheckNameConventionAPI(solution string) []NameConventionCheckResult {
	rt := []NameConventionCheckResult{}
	files := c.GetSolutionFileAttrs(solution)
	for _, file := range files {
		if !file.IsFolder {
			correct, correcetName := ensureFileNameConvention(file.FilePath)
			row := NameConventionCheckResult{}
			row.Correct = correct
			row.CorrectName = correcetName
			row.IncludePath = file.FilePath
			row.File = file
			rt = append(rt, row)
		}
	}

	return rt
}

// CheckNameConvention of the solution
func (c *PDIClient) CheckNameConvention(solution string) {
	count := 0

	for _, r := range c.CheckNameConventionAPI(solution) {
		if !r.Correct {
			count = count + 1
			log.Printf("The name should be %s of file %s\n", r.CorrectName, shortenPath2(r.IncludePath))
		}
	}

	c.exitCode = count
}

// CheckNameConventionToFile output
func (c *PDIClient) CheckNameConventionToFile(solution, output string) {

	tableData := [][]string{}

	for _, r := range c.CheckNameConventionAPI(solution) {
		row := []string{shortenPath2(r.IncludePath), strconv.FormatBool(r.Correct), r.CorrectName}
		tableData = append(tableData, row)

	}

	ss := spreadsheet.New()

	addSheetTo(ss, "Name Convension Check Result", []string{"File", "Correct", "Correct Name"}, tableData)

	ss.SaveToFile(output)

}
