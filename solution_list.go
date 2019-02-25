package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"

	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
)

// GetSolutionsAPI list
func (c *PDIClient) GetSolutionsAPI() []Solution {
	rt := []Solution{}
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

	for _, solution := range solutions {
		solutionInfo := Solution{}

		product := solution.Get("PRODUCT")
		detail := solution.Get("VERSIONS.0")

		solutionInfo.Name = product.Get("PRODUCT").String()
		solutionInfo.Customer = product.Get("PARTNER").String()
		solutionInfo.Contact = product.Get("CONTACT_PERSON").String()
		solutionInfo.Email = product.Get("EMAIL").String()

		solutionInfo.Status = detail.Get("PV_OVERALL_STATUS").String()
		solutionInfo.PatchSolution = detail.Get("PV_1O_PATCH_SOLUTION").String() == "P"
		solutionInfo.Description = detail.Get("PRODUCT_VERSION_TEXTS.0.DDTEXT").String()

		rt = append(rt, solutionInfo)
	}

	return rt
}

func (c *PDIClient) GetSolutionByIDOrDescription(input string) Solution {

	solutions := c.GetSolutionsAPI()
	matched := []Solution{}

	for _, s := range solutions {
		if s.Name == input || s.Description == input {
			matched = append(matched, s)
		}
	}

	switch len(matched) {
	case 0:
		panic(fmt.Errorf("Not found solution with id or description: %s", input))
	case 1:
		return matched[0]
	case 2:
		for _, s := range matched {
			if s.PatchSolution {
				log.Printf(
					"Get solution %s by description %s, default use patch solution",
					s.Name,
					input,
				)
				return s
			}
		}
		panic(fmt.Errorf("Un-expected error"))
	default:
		panic(fmt.Errorf("Un-expected error"))
	}

}

// GetSolutionIDByString for ensure solution ID
func (c *PDIClient) GetSolutionIDByString(input string) string {
	return c.GetSolutionByIDOrDescription(input).Name
}

// ListSolutions detail information
func (c *PDIClient) ListSolutions() *PDIClient {
	solutions := c.GetSolutionsAPI()

	solutionTable := [][]string{}
	for _, solution := range solutions {
		row := []string{solution.Name, solution.Description, strconv.FormatBool(solution.PatchSolution), solution.Status, solution.Customer, solution.Contact, solution.Email}
		solutionTable = append(solutionTable, row)
	}

	// > output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "Description", "Patch", "Status", "Customer", "Contact", "Email"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(solutionTable)
	table.Render()

	return c
}

var commandSolutionList = cli.Command{
	Name:  "list",
	Usage: "list all solutions",
	Action: PDIAction(func(pdiClient *PDIClient, context *cli.Context) {
		pdiClient.ListSolutions()
	}),
}
