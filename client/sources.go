package client

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

type XrepDownloadTask struct {
	xrepPath  string
	localPath string
}

func (c *PDIClient) downloadAndWrite() {

}

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
func (c *PDIClient) DownloadAllSourceTo(solutionName, targetPath string) {
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
	for _, property := range project.PropertyGroup {
		if property.ProjectSourceFolderinXRep != "" {
			xrepPrefix = property.ProjectSourceFolderinXRep
		}
	}
	for _, group := range project.ItemGroup {
		for _, content := range group.Content {
			xrepPath := strings.Replace(filepath.Join(xrepPrefix, content.Include), "\\", "/", -1)
			localPath := strings.Replace(filepath.Join(output, content.Include), "\\", "/", -1)
			downloadList = append(downloadList, XrepDownloadTask{xrepPath, localPath})
		}
	}
	log.Printf("Will download %d files to %s\n", len(downloadList), output)
	// > request and download
	asyncResponses := make([]chan bool, len(downloadList))
	parallexController := make(chan bool, 20)
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
			done <- true
			<-parallexController
		}(task, asyncResponses[idx])
	}

	for _, response := range asyncResponses {
		<-response // ensure all goroutines finished
	}

	log.Println("Done")
}
