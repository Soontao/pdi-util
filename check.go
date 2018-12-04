package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
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

var contentTypeMapping = map[string]string{
	".absl": "ABSL",
	".bo":   "BUSINESS_OBJECT",
	".qry":  "QUERYDEF",
	".xbo":  "EXTENSION_ENTITY",
	".bco":  "BCO",
	".bcc":  "BCSET",
}

// CheckMessage is backend check result
type CheckMessage struct {
	Column   string
	Row      string
	Severity string
	FileName string
	Message  string
}

func (c *PDIClient) backendCheck(xrepPath string) (bool, *[]CheckMessage) {
	canCheck, msgLst := false, []CheckMessage{}

	contentType := contentTypeMapping[filepath.Ext(xrepPath)]

	if contentType != "" {
		canCheck = true
		url := c.xrepPath()
		query := c.query("00163E0115B01DDFB194EC88B8EE8C9B")
		payload := map[string]interface{}{
			"IMPORTING": map[string]interface{}{
				"IV_CONTENT_TYPE": contentType,
				"IT_XREP_PATH":    []string{xrepPath},
			},
		}
		resp, err := req.Post(url, req.BodyJSON(payload), query)
		if err != nil {
			panic(nil)
		}
		respBody, _ := resp.ToString()
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

// CheckBackendMessage information
func (c *PDIClient) CheckBackendMessage(solution string, concurrent int) {
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

	for _, r := range responses {
		_, filename := filepath.Split(r.FileName)
		log.Printf("[%s] %s(%s,%s): %s\n", r.Severity, filename, r.Row, r.Column, r.Message)
	}

	log.Println("Finished")

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
