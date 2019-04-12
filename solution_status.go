package pdiutil

import (
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type SolutionStatus string
type SolutionPhase string

// S_STATUS_IN_DEV Solution In Development
const S_STATUS_IN_DEV = SolutionStatus("1")

// S_STATUS_IN_DEV Assembled
const S_STATUS_ASSEMBLED = SolutionStatus("2")

// S_STATUS_IN_DEPLOYEMENT In Deployement
const S_STATUS_IN_DEPLOYEMENT = SolutionStatus("3")

// S_STATUS_DEPLOYED Deployed
const S_STATUS_DEPLOYED = SolutionStatus("4")

// S_PHASE_ACTIVATION
const S_PHASE_ACTIVATION = SolutionPhase("ACT")

// S_PHASE_DEVELOPMENT
const S_PHASE_DEVELOPMENT = SolutionPhase("DEV")

// SolutionHeader information
type SolutionHeader struct {
	ChangeDateTime  time.Time
	SolutionID      string
	SolutionName    string
	Version         int64
	Status          SolutionStatus
	StatusText      string
	Phase           string
	CanActivation   bool
	CanAssemble     bool
	CanDownload     bool
	NeedCreatePatch bool
	IsRunningJob    bool
	IsCreatingPatch bool
	// Is Solution Enabled
	Enabled bool
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
	enabled := solutionHeader.Get("IS_ENABLED").String() == "X"
	canActivation := solutionHeader.Get("EV_ACT_STATUS").String() == "X"
	canAssemble := solutionHeader.Get("EV_ASSEMBLE_STATUS").String() == "X"
	canDownload := solutionHeader.Get("EV_DOWNLOAD_STATUS").String() == "X"
	isRunningJob := solutionHeader.Get("EV_IS_SPLIT_JOB_RUNNING").String() == "X"
	isCreatingPatch := solutionHeader.Get("EV_IS_PATCH_JOB_RUNNING").String() == "X"
	needCreatePatch := solutionHeader.Get("IS_PATCHSOL_REQUIRED").String() == "X"

	return &SolutionHeader{
		ChangeDateTime:  changeDateTime,
		SolutionID:      solutionID,
		SolutionName:    solutionName,
		Version:         version,
		Status:          SolutionStatus(status),
		StatusText:      statusText,
		Phase:           phase,
		Enabled:         enabled,
		CanActivation:   canActivation,
		CanAssemble:     canAssemble,
		CanDownload:     canDownload,
		IsRunningJob:    isRunningJob,
		NeedCreatePatch: needCreatePatch,
		IsCreatingPatch: isCreatingPatch,
	}

}
