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

// GetSolutionFileList from vs project file
func (c *PDIClient) GetSolutionFileList(solutionName string) *Project {
	url := c.xrepPath()
	query := c.query("00163E0115B01DDFB194EC88B8EDEC9B")
	solutionFilePath := fmt.Sprintf("/%sMAIN/SRC/%s.myproj", solutionName, TrimSuffix(solutionName, "_"))
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
	project := c.GetSolutionFileList(solutionName)
	info := [][]string{}

	for _, group := range project.ItemGroup {
		// ignore folder & bcset
		for _, content := range group.Content {
			row := []string{content.Include}
			info = append(info, row)
		}
	}

	// > output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(info)
	table.Render()
	return c
}

var commandSolution = cli.Command{
	Name:  "solution",
	Usage: "solution related operations",
	Subcommands: []cli.Command{
		commandSolutionList,
		{
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
		},
	},
}
