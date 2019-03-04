package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
)

// Solution information
type Solution struct {
	Name          string
	Description   string
	PatchSolution bool
	Status        string
	Customer      string
	Contact       string
	Email         string
}

// TrimSuffix string
func TrimSuffix(s, suffix string) string {
	rt := s
	if strings.HasSuffix(s, suffix) {
		rt = s[:len(s)-len(suffix)]
	}
	return rt
}

// solution project data cache
var solutionCache = map[string]*Project{}

// GetSolutionFileList from vs project file
// with cache
func (c *PDIClient) GetSolutionFileList(solutionName string) *Project {
	solutionID := c.GetSolutionByIDOrDescription(solutionName).Name
	// with cache
	if p, ok := solutionCache[solutionName]; ok {
		return p
	}
	// else
	url := c.xrepPath()
	query := c.query("00163E0115B01DDFB194EC88B8EDEC9B")
	solutionFilePath := fmt.Sprintf(
		"/%sMAIN/SRC/%s.myproj",
		solutionID,
		TrimSuffix(solutionID, "_"),
	)
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IV_WITH_CONTENT": "X",
			"IV_VIRTUAL_VIEW": "X",
			"IV_PATH":         solutionFilePath,
		},
	}
	resp, err := req.Post(url, req.BodyJSON(payload), query)
	if err != nil {
		panic(nil)
	}
	respBody, _ := resp.ToString()
	if success := gjson.Get(respBody, "EXPORTING.EV_SUCCESS").String(); success != "X" {
		panic(fmt.Sprintf("Not fount project : %s", solutionName))
	}
	projectFileBase64 := gjson.Get(respBody, "EXPORTING.EV_CONTENT").String()
	projectContent, err := base64.StdEncoding.DecodeString(projectFileBase64)
	if err != nil {
		panic(err)
	}
	project := &Project{}
	if err = xml.Unmarshal(projectContent, project); err != nil {
		panic(err)
	}
	return project

}

// ListSolutionAllFiles names
func (c *PDIClient) ListSolutionAllFiles(solutionName string) *PDIClient {
	files := c.GetSolutionXrepFileList(solutionName)
	info := [][]string{}

	for _, f := range files {
		row := []string{f}
		info = append(info, row)
	}

	// > output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(info)
	table.Render()
	return c
}

var commandListSolutionFiles = cli.Command{
	Name:  "files",
	Usage: "list all files in a solution",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "solution, s",
			EnvVar: "SOLUTION_NAME",
			Usage:  "The PDI Solution Name",
		},
	},
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		solutionName := pdiClient.GetSolutionIDByString(context.String("solution"))
		pdiClient.ListSolutionAllFiles(solutionName)
	}),
}

var commandSolution = cli.Command{
	Name:  "solution",
	Usage: "solution related operations",
	Subcommands: []cli.Command{
		commandSolutionList,
		commandListSolutionFiles,
	},
}
