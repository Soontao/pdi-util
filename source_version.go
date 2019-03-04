package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
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

func (v *FileVersion) GetModificationDateTime() time.Time {
	sDateTime := strings.Split(v.Timestamp, ".")[0]
	t, _ := time.Parse("20060102150405", sDateTime)
	return t
}

func (v *FileVersion) GetUserName() string {
	return v.client.GetAUserIDNameByTechID(v.User)
}

func (v *FileVersion) GetVersionContent() *XrepFile {
	return v.client.DownloadVersionFileSource(v.FilePath, *v)
}

// DownloadVersionFileSource will return the remote file content with version information
func (c *PDIClient) DownloadVersionFileSource(xrepPath string, version FileVersion) *XrepFile {

	url := c.xrepPath()
	query := c.query("00163E0115011DDFAEE8C7ADCF082648")
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IS_VERSION_ID": map[string]interface{}{
				"BRANCH":    version.Branch,
				"SOLUTION":  version.Solution,
				"TIMESTAMP": version.Timestamp,
			},
			"IV_VIRTUAL_VIEW": "X",
			"IV_FILE_PATH":    xrepPath,
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
	return &XrepFile{xrepPath, fileContent, attrs}
}

func (c *PDIClient) ViewFileVerionContent(xrepPath string, version FileVersion) {
	file := c.DownloadVersionFileSource(xrepPath, version)
	fmt.Print(string(file.Source))
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

var commandViewFileVersion = cli.Command{
	Name:  "versionview",
	Usage: "view file version content",
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
			Name:   "version, v",
			EnvVar: "VERSION",
			Usage:  "File version string",
		},
	},
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		filename := context.String("filename")
		sVersion := context.String("version")
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))

		if path, found := pdiClient.GetXrepPathByFuzzyName(solutionName, filename); found {
			if version, foundVersion := pdiClient.GetVersionByFuzzyVersion(path, sVersion); foundVersion {
				pdiClient.ViewFileVerionContent(path, version)
			}
		}
	}),
}

var commandListFileVersion = cli.Command{
	Name:  "version",
	Usage: "list file all versions",
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
	},
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		filename := context.String("filename")
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))

		if path, found := pdiClient.GetXrepPathByFuzzyName(solutionName, filename); found {
			pdiClient.ListFileVersions(path)
		}

	}),
}
