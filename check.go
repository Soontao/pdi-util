package main

import (
	"fmt"
	"strconv"

	"baliance.com/gooxml/spreadsheet"
	"github.com/urfave/cli"
)

var commandCheckAll = cli.Command{
	Name:  "all",
	Usage: "do all check to file",
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
			Value:  "check_all.xlsx",
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name",
		},
	},
	Action: PDIAction(func(c *PDIClient, context *cli.Context) {
		solution := context.String("solution")
		concurrent := context.Int("concurrent")
		output := context.String("fileoutput")

		ss := spreadsheet.New()

		backendCheckResponse := c.CheckBackendMessageAPI(solution, concurrent)

		backendTableData := [][]string{}

		translationResponses := c.CheckTranslationAPI(solution, concurrent)

		translationTableData := [][]string{}

		copyRightresponses := c.CheckSolutionCopyrightHeaderAPI(solution, concurrent)

		copyRightTableData := [][]string{}

		nameConventionResponse := c.CheckNameConventionAPI(solution)

		nameConventionTableData := [][]string{}

		for _, r := range nameConventionResponse {
			row := []string{shortenPath2(r.IncludePath), strconv.FormatBool(r.Correct), r.CorrectName}
			nameConventionTableData = append(nameConventionTableData, row)

		}

		for _, r := range backendCheckResponse {

			row := []string{r.GetMessageLevel(), r.GetMessageCategory(), fmt.Sprintf("%s (%s,%s)", shortenPath2(r.FileName), r.Row, r.Column), r.Message}

			backendTableData = append(backendTableData, row)

		}

		for _, r := range translationResponses {

			row := []string{shortenPath2(r.FileName), r.AllTextCount, r.Info["Chinese"].TranslatedCount, r.Info["English"].TranslatedCount}

			translationTableData = append(translationTableData, row)

		}

		for _, r := range copyRightresponses {

			row := []string{shortenPath2(r.File.XrepPath), strconv.FormatBool(r.HaveHeader), r.File.Attributes["~CREATED_BY"], r.File.Attributes["~LAST_CHANGED_BY"]}

			copyRightTableData = append(copyRightTableData, row)

		}

		addSheetTo(ss, "Backend Check Result", []string{"Level", "Category", "Location", "Message"}, backendTableData)

		addSheetTo(ss, "Translation Check Result", []string{"File", "All Field", "Chinese", "English"}, translationTableData)

		addSheetTo(ss, "Copyright Header Check Result", []string{"File", "With Copyright header", "Created By", "Changed By"}, copyRightTableData)

		addSheetTo(ss, "Name Convension Check Result", []string{"File", "Correct", "Correct Name"}, nameConventionTableData)

		ss.SaveToFile(output)

	}),
}

// empty index
// content splited

var commandCheck = cli.Command{
	Name:  "check",
	Usage: "static code check",
	Subcommands: []cli.Command{
		commandCheckAll,
		commandCheckCopyright,
		commandCheckBackend,
		commandCheckTranslation,
		commandCheckNameConvention,
	},
}
