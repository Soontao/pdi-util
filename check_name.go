package main

import (
	"log"
	"strconv"
	"strings"

	"baliance.com/gooxml/spreadsheet"
	"github.com/urfave/cli"
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

// NameConventionCheckResult type
type NameConventionCheckResult struct {
	Correct     bool
	IncludePath string
	CorrectName string
}

// CheckNameConventionAPI of the solution
func (c *PDIClient) CheckNameConventionAPI(solution string) []NameConventionCheckResult {
	rt := []NameConventionCheckResult{}
	project := c.GetSolutionFileList(solution)
	for _, group := range project.ItemGroup {
		for _, bcset := range group.BCSet {
			includePath := strings.SplitN(bcset.Include, "\\", 3)[2]
			correct, correcetName := ensureFileNameConvention(includePath)
			row := NameConventionCheckResult{}
			row.Correct = correct
			row.CorrectName = correcetName
			row.IncludePath = includePath
			rt = append(rt, row)
		}
		for _, content := range group.Content {
			includePath := content.Include
			correct, correcetName := ensureFileNameConvention(includePath)
			row := NameConventionCheckResult{}
			row.Correct = correct
			row.CorrectName = correcetName
			row.IncludePath = includePath
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
	log.Printf("name convension error count: %d", count)
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

var commandCheckNameConvention = cli.Command{
	Name:      "name",
	Usage:     "check name convension",
	UsageText: "check the name convension of source code",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
		cli.StringFlag{
			Name:   "fileoutput, f",
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name",
		},
	},
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		solutionName := context.String("solution")
		output := context.String("fileoutput")
		if output == "" {
			pdiClient.CheckNameConvention(solutionName)
		} else {
			pdiClient.CheckNameConventionToFile(solutionName, output)
		}
	}),
}
