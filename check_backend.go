package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"baliance.com/gooxml/measurement"

	"github.com/urfave/cli"

	"baliance.com/gooxml/spreadsheet"
	"github.com/tidwall/gjson"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var contentTypeMapping = map[string]string{
	".absl":        "ABSL",
	".bo":          "BUSINESS_OBJECT",
	".qry":         "QUERYDEF",
	".xbo":         "EXTENSION_ENTITY",
	".bco":         "BCO",
	".bcc":         "BCSET",
	".uicomponent": "UICOMPONENT",
}

var checkMessageCategoryReg = map[string]*regexp.Regexp{
	"Type Assignment Error":       regexp.MustCompile("Assignment of the type '.*?' to the type '.*?' is not possible."),
	"String Exceed Warning":       regexp.MustCompile("Please ensure that the value doesn’t exceed [\\d]+ characters, otherwise it will be cut off at runtime."),
	"Query Not Supported Warning": regexp.MustCompile("Use of data type '.*?' is not supported in queries"),
	"Type Recommendation":         regexp.MustCompile("Do not store external document data in unrestricted data type '.*?'. Recommendation .*?$"),
	"Performation Warning":        regexp.MustCompile("Performance Alert: .*?$"),
	"Database Length Warning":     regexp.MustCompile("Length of data type '.*?' is restricted to [\\d]+ characters in the data base"),
}

// CheckMessage is backend check result
type CheckMessage struct {
	Column   string
	Row      string
	Severity string
	FileName string
	Message  string
}

// GetMessageCategory string
func (m CheckMessage) GetMessageCategory() string {
	rt := "Unknown"
	for k, r := range checkMessageCategoryReg {
		if r.MatchString(m.Message) {
			return k
		}
	}
	return rt
}

// GetMessageLevel formatted level
// Warning or Error
func (m CheckMessage) GetMessageLevel() string {
	rt := "Unknown"
	switch m.Severity {
	case "W":
		rt = "Warning"
	case "E":
		rt = "Error"
	}
	return rt
}

func (c *PDIClient) backendCheck(xrepPath string) (bool, *[]CheckMessage) {
	canCheck, msgLst := false, []CheckMessage{}

	contentType := contentTypeMapping[filepath.Ext(xrepPath)]

	if contentType != "" {
		canCheck = true
		payload := map[string]interface{}{
			"IMPORTING": map[string]interface{}{
				"IV_CONTENT_TYPE": contentType,
				"IT_XREP_PATH":    []string{xrepPath},
			},
		}

		respBody := c.xrepRequest("00163E0115B01DDFB194EC88B8EE8C9B", payload)
		msgList := gjson.Get(respBody, "EXPORTING.ET_MSG_LIST").Array()
		for _, msg := range msgList {
			checkMessage := CheckMessage{
				Column:   strings.TrimSpace(msg.Get("COLUMN_NO").String()),
				Row:      strings.TrimSpace(msg.Get("LINE_NO").String()),
				Severity: msg.Get("SEVERITY").String(),
				Message:  msg.Get("TEXT").String(),
				FileName: msg.Get("FILE_NAME").String(),
			}
			msgLst = append(msgLst, checkMessage)
		}
	}

	return canCheck, &msgLst
}

// CheckBackendMessageAPI information
func (c *PDIClient) CheckBackendMessageAPI(solution string, concurrent int) []CheckMessage {
	files := c.GetSolutionXrepFileList(solution)
	fileCount := len(files)

	responses := []CheckMessage{}

	asyncResponses := make([]chan *[]CheckMessage, fileCount)
	parallexController := make(chan bool, concurrent)

	bar := pb.New(fileCount)
	bar.ShowBar = false
	bar.Start()
	for idx, file := range files {
		asyncResponses[idx] = make(chan *[]CheckMessage, 1)
		parallexController <- true
		go func(file string, done chan *[]CheckMessage) {
			_, checkMessage := c.backendCheck(file)
			done <- checkMessage
			<-parallexController
			bar.Increment()
		}(file, asyncResponses[idx])
	}

	for _, response := range asyncResponses {
		r := *(<-response)
		responses = append(responses, r...)
	}
	bar.Finish()

	return responses

}

// CheckBackendMessageToFile to output result file
func (c *PDIClient) CheckBackendMessageToFile(solution string, concurrent int, output string) {

	responses := c.CheckBackendMessageAPI(solution, concurrent)

	ss := spreadsheet.New()

	sheet := ss.AddSheet()

	fs := ss.StyleSheet.AddFont()
	fs.SetBold(true)
	cs := ss.StyleSheet.AddCellStyle()
	cs.SetFont(fs)

	sheet.SetName("Backend Check Result")

	header := sheet.AddRow()

	headerSs := []string{"Level", "Category", "Location", "Message"}

	for _, h := range headerSs {
		c := header.AddCell()
		c.SetString(h)
		c.SetStyle(cs)
	}

	for _, r := range responses {
		if r.Severity == "E" {
			// any error will cause exit not as zero
			c.exitCode = c.exitCode + 1
		}
		row := sheet.AddRow()

		row.SetHeightAuto()
		row.AddCell().SetString(r.GetMessageLevel())
		row.AddCell().SetString(r.GetMessageCategory())
		row.AddCell().SetString(fmt.Sprintf("%s (%s,%s)", shortenPath2(r.FileName), r.Row, r.Column))
		row.AddCell().SetString(r.Message)

	}

	sheet.Column(2).SetWidth(measurement.Pixel72 * 350)
	sheet.Column(3).SetWidth(measurement.Pixel72 * 650)
	sheet.Column(4).SetWidth(measurement.Pixel72 * 1370)
	ss.SaveToFile(output)
}

// CheckBackendMessage information
func (c *PDIClient) CheckBackendMessage(solution string, concurrent int) {

	responses := c.CheckBackendMessageAPI(solution, concurrent)

	for _, r := range responses {
		_, filename := filepath.Split(r.FileName)

		if r.Severity == "E" {
			// any error will cause exit not as zero
			c.exitCode = c.exitCode + 1
		}

		log.Printf("[%s]\t%s(%s,%s): %s\n", r.GetMessageLevel(), filename, r.Row, r.Column, r.Message)
	}

}

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
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		solutionName := context.String("solution")
		concurrent := context.Int("concurrent")
		output := context.String("fileoutput")
		if output == "" {
			pdiClient.CheckBackendMessage(solutionName, concurrent)
		} else {
			pdiClient.CheckBackendMessageToFile(solutionName, concurrent, output)
		}
	}),
}
