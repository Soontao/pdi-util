package client

import (
	"os"

	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
)

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
	table.SetHeader([]string{"ID", "Description", "Status", "Customer", "Contact", "Email"})
	table.AppendBulk(solutionTable)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()

	return c
}
