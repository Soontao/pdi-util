package pdiutil

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
)

// CheckBacOutOfDate func
// if the BAC file out of date, assemble will be failed
func (c *PDIClient) CheckBacOutOfDate(solution string) (isOutOfDate bool, errors []error) {
	isOutOfDate = false

	boInformation := c.QueryBOInformation(solution)
	bacInformation := c.GetSolutionBACFile(solution)

	bacMap := map[string]interface{}{}

	for _, c := range bacInformation.Content {
		bacMap[c.ID] = c
	}

	for _, i := range boInformation.Values.RESULTSET.Item {
		// item is BO
		if i.OBJECTCATEGORY != "1" {
			continue
		}

		// like Y7NLADCSY_TE0102E9D042E2B0E15F
		if _, found := bacMap[i.PROXYNAME]; !found {
			// not found in bac file
			e := fmt.Errorf("Not found BO %v in BAC file", i.NAME)
			errors = append(errors, e)
			isOutOfDate = true
		}
	}

	haveBcc := false

	for _, f := range c.GetSolutionXrepFileList(solution) {
		if strings.HasSuffix(f, ".bcc") {
			haveBcc = true
			fNameWithoutExt := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
			if _, found := bacMap[fNameWithoutExt]; !found {
				e := fmt.Errorf("Not found BC Set: %v in BAC file", fNameWithoutExt)
				errors = append(errors, e)
				isOutOfDate = true
			}
		}
	}

	if bacInformation.Bac.VisibleFineTuning == "true" && !haveBcc {
		e := fmt.Errorf("If user selected 'Visible in Fine Tuning', the solution must have at least one BC View")
		errors = append(errors, e)
		isOutOfDate = true
	}

	return isOutOfDate, errors
}

// GetSolutionBACFile content
func (c *PDIClient) GetSolutionBACFile(solution string) *BCPartnerSolution {

	// NEED do more update
	if c.sessionID == "" {
		panic("session id lost")
	}

	files := c.GetSolutionXrepFileList(solution)
	bacFile := ""

	for _, file := range files {
		if strings.HasSuffix(file, ".bac") {
			bacFile = file
			break
		}
	}

	if bacFile == "" {
		panic(fmt.Sprintf("Not found BAC file for solution: %v", solution))
	}

	f := c.DownloadFileSource(bacFile)

	rt := &BCPartnerSolution{}

	if err := xml.Unmarshal(f.Source, rt); err != nil {
		panic(err)
	}

	return rt

}
