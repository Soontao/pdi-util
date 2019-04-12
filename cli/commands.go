package main

import (
	"fmt"
	"log"
	"strconv"

	"baliance.com/gooxml/spreadsheet"
	pdiutil "github.com/Soontao/pdi-util"
	"github.com/urfave/cli"
)

var commandCheckBackend = cli.Command{
	Name:  "backend",
	Usage: "do backend check",
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
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name",
		},
	},
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		concurrent := context.Int("concurrent")
		output := context.String("fileoutput")
		if output == "" {
			pdiClient.CheckBackendMessage(solutionName, concurrent)
		} else {
			pdiClient.CheckBackendMessageToFile(solutionName, concurrent, output)
		}
	}),
}

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
	Action: PDIAction(func(c *pdiutil.PDIClient, context *cli.Context) {
		solution := c.GetSolutionIDByString(context.String("solution"))
		concurrent := context.Int("concurrent")
		output := context.String("fileoutput")

		overviewItems := []pdiutil.OverviewItem{}

		ss := spreadsheet.New()

		backendTableData := [][]string{}

		translationTableData := [][]string{}

		copyRightTableData := [][]string{}

		nameConventionTableData := [][]string{}

		inActiveFilesTableData := [][]string{}

		log.Printf("Starting Backend Check...")

		backendCheckResponse := c.CheckBackendMessageAPI(solution, concurrent)
		backendStatus := pdiutil.OverviewItem{ItemName: "Backend Check Status", ItemDescription: "", ItemStatus: pdiutil.Warning}

		log.Printf("Backend Check Finished")

		log.Printf("Starting Translation Check...")

		translationResponses := c.CheckTranslationAPI(solution, concurrent)

		log.Printf("Translation Check Finished")

		copyRightresponses := c.CheckSolutionCopyrightHeaderAPI(solution, concurrent)
		copyRightStatus := pdiutil.OverviewItem{ItemName: "Copyright Header Check Status", ItemDescription: "", ItemStatus: pdiutil.Warning}

		log.Printf("Copyright Header Check Finished")

		nameConventionResponse := c.CheckNameConventionAPI(solution)
		nameConventionStatus := pdiutil.OverviewItem{ItemName: "Name Convension Check Status", ItemDescription: "All files name correct.", ItemStatus: pdiutil.Successful}

		log.Printf("Name Convention Check Finished")

		inActiveFilesResponse := c.CheckInActiveFilesAPI(solution)
		inActiveFilesStatus := pdiutil.OverviewItem{ItemName: "In-Active Check Status", ItemDescription: "", ItemStatus: pdiutil.Warning}

		log.Printf("In-Active File Check Finished")

		log.Printf("Starting Generating Excel File...")

		// > format data

		// >> name convension
		wrongNameCount := 0

		for _, r := range nameConventionResponse {

			if !r.Correct {
				wrongNameCount++
			}

			row := []string{pdiutil.ShortenPath2(r.IncludePath), strconv.FormatBool(r.Correct), r.CorrectName}
			nameConventionTableData = append(nameConventionTableData, row)
		}

		if wrongNameCount > 0 {
			nameConventionStatus.ItemStatus = pdiutil.Warning
		} else {
			nameConventionStatus.ItemStatus = pdiutil.Successful
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
			case "pdiutil.Warning":
				warningCount++
			}
			row := []string{r.GetMessageLevel(), r.GetMessageCategory(), fmt.Sprintf("%s (%s,%s)", pdiutil.ShortenPath2(r.FileName), r.Row, r.Column), r.Message}
			backendTableData = append(backendTableData, row)
		}
		if errorCount > 0 {
			backendStatus.ItemStatus = pdiutil.FatalError
		} else {
			backendStatus.ItemStatus = pdiutil.Successful
		}

		backendStatus.ItemDescription = fmt.Sprintf("%d errors and %d warnings existed in solution", errorCount, warningCount)

		overviewItems = append(overviewItems, backendStatus)

		// >> translation

		for _, r := range translationResponses {
			row := []string{pdiutil.ShortenPath2(r.FileName), r.AllTextCount, r.Info["Chinese"].TranslatedCount, r.Info["English"].TranslatedCount}
			translationTableData = append(translationTableData, row)
		}

		// >> copyright header

		withoutHeaderCount := 0

		for _, r := range copyRightresponses {
			if !r.HaveHeader {
				withoutHeaderCount++
			}
			row := []string{pdiutil.ShortenPath2(r.File.FilePath), strconv.FormatBool(r.HaveHeader), r.File.CreatedBy, r.File.LastChangedBy}
			copyRightTableData = append(copyRightTableData, row)
		}

		if withoutHeaderCount > 0 {
			copyRightStatus.ItemStatus = pdiutil.Warning
		}
		copyRightStatus.ItemDescription = fmt.Sprintf("%d file lost copyright header", withoutHeaderCount)

		overviewItems = append(overviewItems, copyRightStatus)

		// >> in-active file

		for _, r := range inActiveFilesResponse {
			row := []string{r.File, pdiutil.ShortenPath2(r.FilePath), r.LastChangedBy, r.LastChangedOn}
			inActiveFilesTableData = append(inActiveFilesTableData, row)
		}

		if len(inActiveFilesResponse) > 0 {
			inActiveFilesStatus.ItemDescription = fmt.Sprintf("Found %d in-active files", len(inActiveFilesResponse))
			inActiveFilesStatus.ItemStatus = pdiutil.FatalError
		} else {
			inActiveFilesStatus.ItemDescription = "No in-active files found."
			inActiveFilesStatus.ItemStatus = pdiutil.Successful
		}

		overviewItems = append(overviewItems, inActiveFilesStatus)

		// > generate table & sheets

		pdiutil.AddOverviewSheetTo(ss, overviewItems, c.GetSolutionByIDOrDescription(solution))

		pdiutil.AddSheetTo(ss, "Backend Check Result", []string{"Level", "Category", "Location", "Message"}, backendTableData)

		pdiutil.AddSheetTo(ss, "Translation Check Result", []string{"File", "All Field", "Chinese", "English"}, translationTableData)

		pdiutil.AddSheetTo(ss, "Copyright Header Check Result", []string{"File", "With Copyright header", "Created By", "Changed By"}, copyRightTableData)

		pdiutil.AddSheetTo(ss, "Name Convension Check Result", []string{"File", "Correct", "Correct Name"}, nameConventionTableData)

		pdiutil.AddSheetTo(ss, "InActive Files", []string{"File Name", "File Path", "Last Changed By", "Last Changed On"}, inActiveFilesTableData)

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

var commandSolutionList = cli.Command{
	Name:  "list",
	Usage: "list all solutions",
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		pdiClient.ListSolutions()
	}),
}

var commandListSolutionFiles = cli.Command{
	Name:  "files",
	Usage: "list all files in a solution",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
	},
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		pdiClient.ListSolutionAllFiles(solutionName)
	}),
}

var commandSolution = cli.Command{
	Name:  "solution",
	Usage: "solution related operations",
	Subcommands: []cli.Command{
		commandSolutionList,
		commandListSolutionFiles,
		commandSolutionStatus,
		commandSolutionStatusWatch,
	},
}

var commandCheckCopyright = cli.Command{
	Name:      "header",
	Usage:     "check copyright header",
	UsageText: "\nmake sure all absl & bo have copyright header with following format:\n\n/*\n\tFunction: make sure all absl & bo have copyright header\n\tAuthor: Theo Sun\n\tCopyright: ?\n*/",
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
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name",
		},
	},
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		concurrent := context.Int("concurrent")
		output := context.String("fileoutput")
		if output == "" {
			pdiClient.CheckSolutionCopyrightHeader(solutionName, concurrent)
		} else {
			pdiClient.CheckSolutionCopyrightHeaderToFile(solutionName, concurrent, output)
		}

	}),
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
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		output := context.String("fileoutput")
		if output == "" {
			pdiClient.CheckNameConvention(solutionName)
		} else {
			pdiClient.CheckNameConventionToFile(solutionName, output)
		}
	}),
}

var commandCheckTranslation = cli.Command{
	Name:  "translation",
	Usage: "do translation check",
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
			Name:   "language, l",
			EnvVar: "LANGUAGE",
			Value:  "Chinese",
			Usage:  "target language to check",
		},
		cli.StringFlag{
			Name:   "fileoutput, f",
			EnvVar: "FILENAME_OUTPUT",
			Usage:  "output file name",
		},
	},
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		concurrent := context.Int("concurrent")
		language := context.String("language")
		output := context.String("fileoutput")
		if output == "" {
			pdiClient.CheckTranslation(solutionName, concurrent, language)
		} else {
			pdiClient.CheckTranslationToFile(solutionName, concurrent, language, output)
		}

	}),
}

var commandDownloadSource = cli.Command{
	Name:  "download",
	Usage: "download all files in a solution",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
		cli.StringFlag{
			Name:   "output, o",
			EnvVar: "OUTPUT",
			Value:  "output",
			Usage:  "Output directory",
		},
		cli.IntFlag{
			Name:   "concurrent, c",
			EnvVar: "DOWNLOAD_CONCURRENT",
			Value:  35,
			Usage:  "concurrent goroutines number",
		},
		cli.BoolFlag{
			Name:   "pretty, f",
			EnvVar: "PRETTY",
			Usage:  "pretty xml files",
		},
	},
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		output := context.String("output")
		concurrent := context.Int("concurrent")
		pretty := context.Bool("pretty")
		pdiClient.DownloadAllSourceTo(solutionName, output, concurrent, pretty)
	}),
}

var commandSource = cli.Command{
	Name:  "source",
	Usage: "source code related operations",
	Subcommands: []cli.Command{
		commandListFileVersion,
		commandDownloadSource,
	},
}

var commandListFileVersion = cli.Command{
	Name:  "version",
	Usage: "list/view/diff file versions",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
		cli.StringFlag{
			Name:   "filename, f",
			EnvVar: "VERSION_FILE_NAME",
			Usage:  "The target file xrep path/file name",
		},
		cli.StringFlag{
			Name:  "from",
			Usage: "Version From",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "Version To",
		},
		cli.StringFlag{
			Name:   "targetversion, v",
			EnvVar: "VERSION",
			Usage:  "File version string",
		},
	},
	Action: PDIAction(func(pdiClient *pdiutil.PDIClient, context *cli.Context) {
		filename := context.String("filename")
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		sVersionFrom := context.String("from")
		sVersion := context.String("targetversion")
		sVersionTo := context.String("to")

		if sVersionFrom != "" && sVersionTo != "" {
			if path, found := pdiClient.GetXrepPathByFuzzyName(solutionName, filename); found {
				if versionFrom, foundVersion := pdiClient.GetVersionByFuzzyVersion(path, sVersionFrom); foundVersion {
					if versionTo, foundVersion := pdiClient.GetVersionByFuzzyVersion(path, sVersionTo); foundVersion {
						pdiClient.DiffFileVersion(versionFrom, versionTo)
					}
				}
			}
		} else if sVersion != "" {

			if path, found := pdiClient.GetXrepPathByFuzzyName(solutionName, filename); found {
				if version, foundVersion := pdiClient.GetVersionByFuzzyVersion(path, sVersion); foundVersion {
					pdiClient.ViewFileVerionContent(version)
				}
			}

		} else {

			if path, found := pdiClient.GetXrepPathByFuzzyName(solutionName, filename); found {
				pdiClient.ListFileVersions(path)
			}

		}

	}),
}
