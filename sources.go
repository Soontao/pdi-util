package pdiutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/spkg/bom"
	"golang.org/x/sync/semaphore"

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

func (f *XrepFile) String() string {
	return string(f.Source)
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
	// remove MS UTF8-BOM bytes
	fileContent = bom.Clean(fileContent)
	return &XrepFile{xrepPath, fileContent, attrs}
}

// GetXrepPathByFuzzyName func
func (c *PDIClient) GetXrepPathByFuzzyName(solution, s string) (xrepPath string, existAndUnique bool) {
	matched := []string{}

	for _, xFilePath := range c.GetSolutionXrepFileList(solution) {
		if strings.Contains(xFilePath, s) {
			matched = append(matched, xFilePath)
		}
	}

	switch len(matched) {
	case 0:
		log.Printf("Not found any file with name: %s", s)
		existAndUnique = false
	case 1:
		xrepPath = matched[0]
		existAndUnique = true
	default:
		log.Println("More than one files matched name: " + s)
		for _, m := range matched {
			log.Println(m)
		}
		existAndUnique = false

	}

	return xrepPath, existAndUnique
}

// GetSolutionXrepFileList in xrep
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

// fetchSources list
func (c *PDIClient) fetchSources(xrepPathes []string, concurrent int) []*XrepFile {

	count := len(xrepPathes)

	rt := make([]*XrepFile, len(xrepPathes))

	if count > 0 { // with something need to be download
		log.Printf("Downloading %v items from repository", count)

		bar := pb.StartNew(len(xrepPathes))

		var wg sync.WaitGroup

		ctx := context.TODO()

		sem := semaphore.NewWeighted(int64(concurrent))

		for idx, xrepPath := range xrepPathes {
			wg.Add(1)
			sem.Acquire(ctx, 1)

			go func(i string, s *semaphore.Weighted, w *sync.WaitGroup, payloadIndex int) {

				defer wg.Done()
				rt[payloadIndex] = c.DownloadFileSource(i)
				s.Release(1)
				bar.Increment()

			}(xrepPath, sem, &wg, idx)
		}

		wg.Wait() // wait all requests finished
		bar.Finish()
		log.Println("Download Finished")
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
	xrepFiles := c.fetchSources(c.GetSolutionXrepFileList(solutionName), concurrent)

	for _, xrepFile := range xrepFiles {
		xrepPath := xrepFile.XrepPath
		localPath := strings.Replace(filepath.Join(output, xrepPath), "\\", "/", -1)

		sourceContent := xrepFile.Source
		os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			f, err := os.Create(localPath)
			if err != nil {
				panic(err)
			}
			f.Close()
		}
		if err := ioutil.WriteFile(localPath, sourceContent, 0644); err != nil {
			panic(err)
		}
	}

}
