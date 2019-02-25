package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// XrepDownloadTask is download & write task
type XrepDownloadTask struct {
	xrepPath  string
	localPath string
}

// XrepFile type
type XrepFile struct {
	XrepPath   string
	Source     []byte
	Attributes map[string]string
}

// ListAllLockedFile list all checked out files name
func (c *PDIClient) ListAllLockedFile(solutionName string) []string {
	rt := []string{}
	// to do
	return rt
}

// DownloadFileSource will return the remote file content
func (c *PDIClient) DownloadFileSource(xrepPath string) *XrepFile {

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
	attrs := map[string]string{}
	attrsList := gjson.Get(respBody, "EXPORTING.ET_ATTR").Array()
	for _, attr := range attrsList {
		attrs[attr.Get("NAME").String()] = attr.Get("VALUE").String()
	}
	base64Content := gjson.Get(respBody, "EXPORTING.EV_CONTENT").String()
	fileContent, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		panic(err)
	}
	return &XrepFile{xrepPath, fileContent, attrs}
}

// GetSolutionFileList in xrep
func (c *PDIClient) GetSolutionXrepFileList(solutionName string) []string {
	rt := []string{}
	project := c.GetSolutionFileList(solutionName)
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
			rt = append(rt, xrepPath)
		}
		// Common Files
		for _, content := range group.Content {
			xrepPath := strings.Replace(filepath.Join(xrepPrefix, content.Include), "\\", "/", -1)
			rt = append(rt, xrepPath)
		}
	}

	return rt
}

// DownloadAllSourceTo directory
func (c *PDIClient) DownloadAllSourceTo(solutionName, targetPath string, concurrent int, pretty bool) {
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
	xrepFiles := c.GetSolutionXrepFileList(solutionName)
	downloadList := []XrepDownloadTask{}

	for _, xrepFile := range xrepFiles {
		xrepPath := xrepFile
		localPath := strings.Replace(filepath.Join(output, xrepFile), "\\", "/", -1)
		downloadList = append(downloadList, XrepDownloadTask{xrepPath, localPath})
	}

	fileCount := len(downloadList)

	log.Printf("Will download %d files to %s\n", fileCount, output)
	// > progress ui support
	// > request and download
	asyncResponses := make([]chan bool, fileCount)
	parallexController := make(chan bool, concurrent)
	for idx, task := range downloadList {
		asyncResponses[idx] = make(chan bool, 1)
		parallexController <- true
		go func(task XrepDownloadTask, done chan bool) {
			source := c.DownloadFileSource(task.xrepPath)
			sourceContent := source.Source
			os.MkdirAll(filepath.Dir(task.localPath), os.ModePerm)
			if _, err := os.Stat(task.localPath); os.IsNotExist(err) {
				f, err := os.Create(task.localPath)
				if err != nil {
					panic(err)
				}
				f.Close()
			}
			if err := ioutil.WriteFile(task.localPath, sourceContent, 0644); err != nil {
				panic(err)
			}
			done <- true
			<-parallexController
		}(task, asyncResponses[idx])
	}

	for _, response := range asyncResponses {
		<-response // ensure all goroutines finished
	}
}

var commandSource = cli.Command{
	Name:  "source",
	Usage: "source code related operations",
	Subcommands: []cli.Command{
		{
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
			Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
				solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
				output := context.String("output")
				concurrent := context.Int("concurrent")
				pretty := context.Bool("pretty")
				pdiClient.DownloadAllSourceTo(solutionName, output, concurrent, pretty)
			}),
		},
	},
}
