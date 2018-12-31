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
		solution := c.GetSolutionIDByString(context.String("solution"))
		concurrent := context.Int("concurrent")
		output := context.String("fileoutput")

		overviewItems := []OverviewItem{}

		ss := spreadsheet.New()

		backendTableData := [][]string{}

		translationTableData := [][]string{}

		copyRightTableData := [][]string{}

		nameConventionTableData := [][]string{}

		inActiveFilesTableData := [][]string{}

		log.Printf("Starting Backend Check...")

		backendCheckResponse := c.CheckBackendMessageAPI(solution, concurrent)
		backendStatus := OverviewItem{"Backend Check Status", "", Warning}

		log.Printf("Backend Check Finished")

		log.Printf("Starting Translation Check...")

		translationResponses := c.CheckTranslationAPI(solution, concurrent)

		log.Printf("Translation Check Finished")

		copyRightresponses := c.CheckSolutionCopyrightHeaderAPI(solution, concurrent)
		copyRightStatus := OverviewItem{"Copyright Header Check Status", "", Warning}

		log.Printf("Copyright Header Check Finished")

		nameConventionResponse := c.CheckNameConventionAPI(solution)
		nameConventionStatus := OverviewItem{"Name Convension Check Status", "All files name correct.", Successful}

		log.Printf("Name Convention Check Finished")

		inActiveFilesResponse := c.CheckInActiveFilesAPI(solution)
		inActiveFilesStatus := OverviewItem{"In-Active Check Status", "", Warning}

		log.Printf("In-Active File Check Finished")

		log.Printf("Starting Generating Excel File...")

		// > format data

		// >> name convension
		wrongNameCount := 0

		for _, r := range nameConventionResponse {

			if !r.Correct {
				wrongNameCount++
			}

			row := []string{shortenPath2(r.IncludePath), strconv.FormatBool(r.Correct), r.CorrectName}
			nameConventionTableData = append(nameConventionTableData, row)
		}

		if wrongNameCount > 0 {
			nameConventionStatus.ItemStatus = Warning
		} else {
			nameConventionStatus.ItemStatus = Successful
		}

		nameConventionStatus.ItemDescription = fmt.Sprintf("%d file filenames are incorrect.", wrongNameCount)

		overviewItems = append(overviewItems, nameConventionStatus)

		// >> backend check

		errorCount := 0
		warningCount := 0

		for _, r := range backendCheckResponse {
			switch r.GetMessageLevel() {
			case "Error":
				errorCount++
			case "Warning":
				warningCount++
			}
			row := []string{r.GetMessageLevel(), r.GetMessageCategory(), fmt.Sprintf("%s (%s,%s)", shortenPath2(r.FileName), r.Row, r.Column), r.Message}
			backendTableData = append(backendTableData, row)
		}
		if errorCount > 0 {
			backendStatus.ItemStatus = FatalError
		} else {
			backendStatus.ItemStatus = Successful
		}

		backendStatus.ItemDescription = fmt.Sprintf("%d errors and %d warnings existed in solution", errorCount, warningCount)

		overviewItems = append(overviewItems, backendStatus)

		// >> translation

		for _, r := range translationResponses {
			row := []string{shortenPath2(r.FileName), r.AllTextCount, r.Info["Chinese"].TranslatedCount, r.Info["English"].TranslatedCount}
			translationTableData = append(translationTableData, row)
		}

		// >> copyright header

		withoutHeaderCount := 0

		for _, r := range copyRightresponses {
			if !r.HaveHeader {
				withoutHeaderCount++
			}
			row := []string{shortenPath2(r.File.XrepPath), strconv.FormatBool(r.HaveHeader), r.File.Attributes["~CREATED_BY"], r.File.Attributes["~LAST_CHANGED_BY"]}
			copyRightTableData = append(copyRightTableData, row)
		}

		if withoutHeaderCount > 0 {
			copyRightStatus.ItemStatus = Warning
		}
		copyRightStatus.ItemDescription = fmt.Sprintf("%d file lost copyright header", withoutHeaderCount)

		overviewItems = append(overviewItems, copyRightStatus)

		// >> in-active file

		for _, r := range inActiveFilesResponse {
			row := []string{r.File, shortenPath2(r.FilePath), r.LastChangedBy, r.LastChangedOn}
			inActiveFilesTableData = append(inActiveFilesTableData, row)
		}

		if len(inActiveFilesResponse) > 0 {
			inActiveFilesStatus.ItemDescription = fmt.Sprintf("Found %d in-active files", len(inActiveFilesResponse))
			inActiveFilesStatus.ItemStatus = FatalError
		} else {
			inActiveFilesStatus.ItemDescription = "No in-active files found."
			inActiveFilesStatus.ItemStatus = Successful
		}

		overviewItems = append(overviewItems, inActiveFilesStatus)

		// > generate table & sheets

		addOverviewSheetTo(ss, overviewItems)

		addSheetTo(ss, "Backend Check Result", []string{"Level", "Category", "Location", "Message"}, backendTableData)

		addSheetTo(ss, "Translation Check Result", []string{"File", "All Field", "Chinese", "English"}, translationTableData)

		addSheetTo(ss, "Copyright Header Check Result", []string{"File", "With Copyright header", "Created By", "Changed By"}, copyRightTableData)

		addSheetTo(ss, "Name Convension Check Result", []string{"File", "Correct", "Correct Name"}, nameConventionTableData)

		addSheetTo(ss, "InActive Files", []string{"File Name", "File Path", "Last Changed By", "Last Changed On"}, inActiveFilesTableData)

		ss.SaveToFile(output)

		log.Printf("Saving Check Result File to %s", output)

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
