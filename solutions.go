package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
)

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

// ListSolutions detail information
func (c *PDIClient) ListSolutions() *PDIClient {
	url := c.xrepPath()
	query := c.query("00163E0115B01DDFB194E54BB7204C9B")
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IV_GET":           "X",
			"IV_PRODUCT_NAME":  nil,
			"IV_SOLUTION_TYPE": "2",
			"IV_USER":          c.ivUser,
		},
	}
	resp, err := req.Post(url, req.BodyJSON(payload), query)
	if err != nil {
		panic(nil)
	}
	respBody, _ := resp.ToString()
	solutions := gjson.Get(respBody, "EXPORTING.ET_PRODUCTS").Array()

	solutionTable := [][]string{}

	for _, solution := range solutions {
		product := solution.Get("PRODUCT")
		detail := solution.Get("VERSIONS.0")

		solutionID := product.Get("PRODUCT").String()
		customer := product.Get("PARTNER").String()
		contact := product.Get("CONTACT_PERSON").String()
		email := product.Get("EMAIL").String()

		solutionStatus := detail.Get("PV_OVERALL_STATUS").String()
		solutionDescription := detail.Get("PRODUCT_VERSION_TEXTS.0.DDTEXT").String()
		row := []string{solutionID, solutionDescription, solutionStatus, customer, contact, email}
		solutionTable = append(solutionTable, row)
	}

	// > output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "Description", "Status", "Customer", "Contact", "Email"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(solutionTable)
	table.Render()

	return c
}
