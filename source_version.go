package pdiutil

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/tidwall/gjson"
)

// FileVersion describe a file version information
type FileVersion struct {
	Solution   string
	Branch     string
	Timestamp  string
	Action     string
	ContentID  string
	User       string
	ObjectType string
	FilePath   string
	client     *PDIClient
}

// GetModificationDateTime obj
func (v *FileVersion) GetModificationDateTime() time.Time {
	return ParseXrepDateString(v.Timestamp)
}

func (v *FileVersion) GetUserName() string {
	return v.client.GetAUserIDNameByTechID(v.User)
}

func (v *FileVersion) GetVersionContent() *XrepFile {
	return v.client.DownloadVersionFileSource(v.FilePath, v.Branch, v.Solution, v.Timestamp)
}

// DownloadVersionFileSource will return the remote file content with version information
func (c *PDIClient) DownloadVersionFileSource(path, branch, solution, timestamp string) *XrepFile {

	url := c.xrepPath()
	query := c.query("00163E0115011DDFAEE8C7ADCF082648")
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IS_VERSION_ID": map[string]interface{}{
				"BRANCH":    branch,
				"SOLUTION":  solution,
				"TIMESTAMP": timestamp,
			},
			"IV_VIRTUAL_VIEW": "X",
			"IV_FILE_PATH":    path,
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
	base64Content := gjson.Get(respBody, "EXPORTING.EV_FILE_CONTENT").String()
	fileContent, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		panic(err)
	}
	return &XrepFile{path, fileContent, attrs}
}

// ViewFileVersionContent text
//
// print it to stdout
func (c *PDIClient) ViewFileVersionContent(version FileVersion) {
	file := version.GetVersionContent()
	fmt.Println(string(file.Source))
}

// ListFileVersionsAPI information
func (c *PDIClient) ListFileVersionsAPI(xrepPath string) []FileVersion {
	rt := []FileVersion{}
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IV_CLIENT":       c.sapClient,
			"IV_FILE_PATH":    xrepPath,
			"IV_USER":         c.ivUser,
			"IV_VIRTUAL_VIEW": "X",
		},
	}
	respBody := c.xrepRequest("00163E0115011DDFAEE8C7ADCF062648", payload)
	historyList := gjson.Get(respBody, "EXPORTING.ET_VERSION_HISTORY").Array()

	for _, h := range historyList {
		rt = append(rt, FileVersion{
			Solution:   h.Get("SOLUTION").String(),
			Branch:     h.Get("BRANCH").String(),
			Timestamp:  h.Get("TIMESTAMP").String(),
			Action:     h.Get("ACTION").String(),
			User:       h.Get("CREATED_BY").String(),
			ContentID:  h.Get("CONT_ID").String(),
			ObjectType: h.Get("TROBJTYPE").String(),
			FilePath:   xrepPath,
			client:     c,
		})
	}

	return rt
}

func (c *PDIClient) GetVersionByFuzzyVersion(xrepPath, sVersion string) (rt FileVersion, existAndUnique bool) {
	matched := []FileVersion{}
	for _, v := range c.ListFileVersionsAPI(xrepPath) {
		if strings.Contains(v.Timestamp, sVersion) {
			matched = append(matched, v)
		}
	}

	switch len(matched) {
	case 0:
		log.Printf("Not found any file with name: %s", sVersion)
		existAndUnique = false
	case 1:
		rt = matched[0]
		existAndUnique = true
	default:
		log.Println("More than one files matched name: " + sVersion)
		for _, m := range matched {
			log.Println(m)
		}
		existAndUnique = false

	}

	return rt, existAndUnique

}

// ListFileVersions to console
func (c *PDIClient) ListFileVersions(xrepPath string) {
	versions := c.ListFileVersionsAPI(xrepPath)
	info := [][]string{}

	for _, v := range versions {
		row := []string{
			v.GetModificationDateTime().String(),
			v.Action,
			v.GetUserName(),
			v.Timestamp,
		}
		info = append(info, row)
	}

	// > output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date Time", "Action", "User", "VersionID"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(info)
	table.Render()
}

// DiffFileVersion
func (c *PDIClient) DiffFileVersion(from, to FileVersion) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(from.GetVersionContent().Source), string(to.GetVersionContent().Source), false)
	fmt.Println(dmp.DiffPrettyText(diffs))
}
