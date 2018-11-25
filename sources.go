package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// XrepDownloadTask is download & write task
type XrepDownloadTask struct {
	xrepPath  string
	localPath string
}

// DownloadFileSource will return the remote file content
func (c *PDIClient) DownloadFileSource(xrepPath string) []byte {

	url := c.xrepPath()
	query := c.query("00163E0115B01DDFB194EC88B8EDEC9B")
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IV_VIRTUAL_VIEW": "X",
			"IV_WITH_CONTENT": "X",
			"IV_PATH":         xrepPath,
		},
	}
	resp, err := req.Post(url, req.BodyJSON(payload), query)
	if err != nil {
		panic(nil)
	}
	respBody, _ := resp.ToString()
	base64Content := gjson.Get(respBody, "EXPORTING.EV_CONTENT").String()
	fileContent, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		panic(err)
	}
	return fileContent
}

// DownloadAllSourceTo directory
func (c *PDIClient) DownloadAllSourceTo(solutionName, targetPath string, concurrent int) {
	// > process output target
	ensure(targetPath, "targetPath")
	output := ""
	pwd, _ := os.Getwd()
	if filepath.IsAbs(targetPath) {
		output = targetPath
	} else {
		output = filepath.Join(pwd, targetPath)
	}
	os.MkdirAll(output, os.ModePerm)
	// > get project file list
	project := c.GetSolutionFileList(solutionName)
	downloadList := []XrepDownloadTask{}
	xrepPrefix := ""
	bcPrefix := ""
	for _, property := range project.PropertyGroup {
		if property.ProjectSourceFolderinXRep != "" {
			xrepPrefix = property.ProjectSourceFolderinXRep
		}
		if property.BCSourceFolderInXRep != "" {
			bcPrefix = property.BCSourceFolderInXRep
		}
	}
	for _, group := range project.ItemGroup {
		// Bussiness Configuration Files
		for _, bc := range group.BCSet {
			realPath := strings.TrimPrefix(bc.Include, fmt.Sprintf("..\\%sBC\\", solutionName))
			xrepPath := strings.Replace(filepath.Join(bcPrefix, realPath), "\\", "/", -1)
			localPath := strings.Replace(filepath.Join(output, realPath), "\\", "/", -1)
			downloadList = append(downloadList, XrepDownloadTask{xrepPath, localPath})
		}
		// Common Files
		for _, content := range group.Content {
			xrepPath := strings.Replace(filepath.Join(xrepPrefix, content.Include), "\\", "/", -1)
			localPath := strings.Replace(filepath.Join(output, content.Include), "\\", "/", -1)
			downloadList = append(downloadList, XrepDownloadTask{xrepPath, localPath})
		}
	}

	fileCount := len(downloadList)
	log.Printf("Will download %d files to %s\n", fileCount, output)
	// > progress ui support
	bar := pb.New(fileCount)
	bar.ShowBar = false
	bar.Start()
	// > request and download
	asyncResponses := make([]chan bool, fileCount)
	parallexController := make(chan bool, concurrent)
	for idx, task := range downloadList {
		asyncResponses[idx] = make(chan bool, 1)
		parallexController <- true
		go func(task XrepDownloadTask, done chan bool) {
			source := c.DownloadFileSource(task.xrepPath)
			os.MkdirAll(filepath.Dir(task.localPath), os.ModePerm)
			if _, err := os.Stat(task.localPath); os.IsNotExist(err) {
				f, err := os.Create(task.localPath)
				if err != nil {
					panic(err)
				}
				f.Close()
			}
			if err := ioutil.WriteFile(task.localPath, source, 0644); err != nil {
				panic(err)
			}
			bar.Increment()
			done <- true
			<-parallexController
		}(task, asyncResponses[idx])
	}

	for _, response := range asyncResponses {
		<-response // ensure all goroutines finished
	}
	bar.Finish()
	log.Println("Done")
}
