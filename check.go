package main

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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
	".absl":        "ABSL",
	".bo":          "BUSINESS_OBJECT",
	".qry":         "QUERYDEF",
	".xbo":         "EXTENSION_ENTITY",
	".bco":         "BCO",
	".bcc":         "BCSET",
	".uicomponent": "UICOMPONENT",
}

// CheckMessage is backend check result
type CheckMessage struct {
	Column   string
	Row      string
	Severity string
	FileName string
	Message  string
}

// GetMessageLevel formatted level
// Warning or Error
func (m CheckMessage) GetMessageLevel() string {
	rt := "UNKNOWN"
	switch m.Severity {
	case "W":
		rt = "Warning"
	case "E":
		rt = "Error"
	}
	return rt
}

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
	bar.Finish()

	return responses

}

// CheckTranslation information
func (c *PDIClient) CheckTranslation(solution string, concurrent int, language string) {

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
			log.Printf("For language %s, translated %d text of %d, file(%s)\n", language, translatedCount, targetCount, filename)
		} else {
			log.Printf("For language %s, full translated, file(%s)\n", language, filename)
		}

	}

	log.Println("Finished")

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

// CheckBackendMessage information
func (c *PDIClient) CheckBackendMessage(solution string, concurrent int) {

	responses := c.CheckBackendMessageAPI(solution, concurrent)

	for _, r := range responses {
		_, filename := filepath.Split(r.FileName)
		log.Printf("[%s]\t%s(%s,%s): %s\n", r.GetMessageLevel(), filename, r.Row, r.Column, r.Message)
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
