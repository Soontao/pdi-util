package pdiutil

import (
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// SolutionHeader information
type SolutionHeader struct {
	ChangeDateTime time.Time
	SolutionID     string
	SolutionName   string
	Version        int64
	Status         string
	StatusText     string
	Phase          string
}

// GetSolutionStatus exported
func (c *PDIClient) GetSolutionStatus(solution string) *SolutionHeader {
	solutionID := c.GetSolutionIDByString(solution)
	payload := JSONObject{
		"IMPORTING": JSONObject{
			"IV_LANGUAGE":      "E",
			"IV_LOGLEVEL":      "0",
			"IV_SOLUTION_NAME": solutionID,
		},
	}

	respBody := c.xrepRequest("00163E0123D21EE092C936CC65A49BA4", payload)
	solutionHeader := gjson.Get(respBody, "EXPORTING.ES_SOLUTION_HEADER")

	changeDateTime := ParseXrepDateString(solutionHeader.Get("CHANGE_DATETIME").String())
	solutionName := solutionHeader.Get("DESCRIPTION").String()
	version, _ := strconv.ParseInt(solutionHeader.Get("VERSION_ID").String(), 10, 64)
	status := solutionHeader.Get("VERSION_STATUS").String()
	statusText := solutionHeader.Get("VERSION_STATUS_TEXT").String()
	phase := solutionHeader.Get("PHASE").String()

	return &SolutionHeader{
		ChangeDateTime: changeDateTime,
		SolutionID:     solutionID,
		SolutionName:   solutionName,
		Version:        version,
		Status:         status,
		StatusText:     statusText,
		Phase:          phase,
	}

}
