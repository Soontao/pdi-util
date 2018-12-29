package main

import (
	"fmt"
	"log"
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

		backendTableData := [][]string{}

		translationTableData := [][]string{}

		copyRightTableData := [][]string{}

		nameConventionTableData := [][]string{}

		inActiveFilesTableData := [][]string{}

		backendCheckResponse := c.CheckBackendMessageAPI(solution, concurrent)

		log.Printf("Backend Check Finished")

		translationResponses := c.CheckTranslationAPI(solution, concurrent)

		log.Printf("Translation Check Finished")

		copyRightresponses := c.CheckSolutionCopyrightHeaderAPI(solution, concurrent)

		log.Printf("Copyright Header Check Finished")

		nameConventionResponse := c.CheckNameConventionAPI(solution)

		log.Printf("Name Convention Check Finished")

		inActiveFilesResponse := c.CheckInActiveFilesAPI(solution)

		log.Printf("In-Active File Check Finished")

		log.Printf("Start Generating Excel File...")

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

		for _, r := range inActiveFilesResponse {
			row := []string{r.File, shortenPath2(r.FilePath), r.LastChangedBy, r.LastChangedOn}
			inActiveFilesTableData = append(inActiveFilesTableData, row)
		}

		addSheetTo(ss, "Backend Check Result", []string{"Level", "Category", "Location", "Message"}, backendTableData)

		addSheetTo(ss, "Translation Check Result", []string{"File", "All Field", "Chinese", "English"}, translationTableData)

		addSheetTo(ss, "Copyright Header Check Result", []string{"File", "With Copyright header", "Created By", "Changed By"}, copyRightTableData)

		addSheetTo(ss, "Name Convension Check Result", []string{"File", "Correct", "Correct Name"}, nameConventionTableData)

		addSheetTo(ss, "InActive Files", []string{"File Name", "File Path", "Last Changed By", "Last Changed On"}, inActiveFilesTableData)

		ss.SaveToFile(output)

		log.Printf("Save Check Result File to %s", output)

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
