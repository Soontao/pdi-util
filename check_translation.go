package main

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"baliance.com/gooxml/spreadsheet"
	"github.com/urfave/cli"

	"github.com/tidwall/gjson"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// TranslationStatus message
type TranslationStatus struct {
	FileName     string
	AllTextCount string
	Info         map[string]TranslationStatusInfo
}

// TranslationStatusInfo detail
type TranslationStatusInfo struct {
	Language        string
	TranslatedCount string
}

var translationCheckList = map[string]bool{
	".uicomponent": true, ".bo": true, ".codelist": true,
}

func (c *PDIClient) translationInformation(xrepPath string) (bool, *TranslationStatus) {
	canCheck := translationCheckList[filepath.Ext(xrepPath)]
	rt := &TranslationStatus{
		FileName: xrepPath,
		Info:     map[string]TranslationStatusInfo{},
	}

	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IT_PATH": []string{xrepPath},
		},
	}

	if canCheck {
		respBody := c.xrepRequest("00163E01138A1EE0AFEA287164321C26", payload)
		textCount := strings.TrimSpace(gjson.Get(respBody, "EXPORTING.EV_NUMBER_OF_TEXTS").String())
		rt.AllTextCount = textCount
		checkInfoList := gjson.Get(respBody, "EXPORTING.ET_CHECK_INFO").Array()
		for _, jsonInfo := range checkInfoList {
			info := TranslationStatusInfo{}
			info.Language = jsonInfo.Get("LANGUAGE").String()
			info.TranslatedCount = strings.TrimSpace(jsonInfo.Get("TEXTCOUNT").String())

			rt.Info[info.Language] = info
		}

	}

	return canCheck, rt
}

// CheckTranslationAPI used for programming
func (c *PDIClient) CheckTranslationAPI(solution string, concurrent int) []TranslationStatus {
	files := c.GetSolutionXrepFileList(solution)
	fileCount := len(files)

	responses := []TranslationStatus{}

	asyncResponses := make([]chan *TranslationStatus, fileCount)
	parallexController := make(chan bool, concurrent)

	bar := pb.New(fileCount)
	bar.ShowBar = false
	bar.Start()
	for idx, file := range files {
		asyncResponses[idx] = make(chan *TranslationStatus, 1)
		parallexController <- true
		go func(file string, done chan *TranslationStatus) {
			canCheck, checkMessage := c.translationInformation(file)
			if canCheck {
				done <- checkMessage
			} else {
				done <- nil
			}
			<-parallexController
			bar.Increment()
		}(file, asyncResponses[idx])
	}

	for _, response := range asyncResponses {
		r := (<-response)
		if r != nil {
			responses = append(responses, *r)
		}
	}
	bar.FinishPrint("Translation check finished")

	return responses

}

// CheckTranslationToFile output
func (c *PDIClient) CheckTranslationToFile(solution string, concurrent int, language, outputFile string) {

	responses := c.CheckTranslationAPI(solution, concurrent)

	tableData := [][]string{}

	for _, r := range responses {

		row := []string{shortenPath2(r.FileName), r.AllTextCount, r.Info["Chinese"].TranslatedCount, r.Info["English"].TranslatedCount}

		tableData = append(tableData, row)

	}

	ss := spreadsheet.New()

	addSheetTo(ss, "Translation Check Result", []string{"File", "All Field", "Chinese", "English"}, tableData)

	ss.SaveToFile(outputFile)

}

// CheckTranslation information
func (c *PDIClient) CheckTranslation(solution string, concurrent int, language string) {

	// need check language here

	responses := c.CheckTranslationAPI(solution, concurrent)

	for _, r := range responses {
		_, filename := filepath.Split(r.FileName)
		targetCount, err := strconv.Atoi(r.AllTextCount)

		if err != nil {
			panic(err)
		}

		translatedInfo := r.Info[language]
		translatedCount, err := strconv.Atoi(translatedInfo.TranslatedCount)
		if err != nil {
			panic(err)
		}
		if translatedCount < targetCount {
			c.exitCode = c.exitCode + 1
			log.Printf("For language %s, %d of %d texts have beed translated, file(%s)\n", language, translatedCount, targetCount, filename)
		} else {
			log.Printf("For language %s, fully translated, file(%s)\n", language, filename)
		}

	}

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
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		solutionName := context.String("solution")
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
