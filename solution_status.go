package pdiutil

import (
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type SolutionStatus string
type SolutionPhase string
type SolutionPhaseStatus string

// S_STATUS_IN_DEV Solution In Development
const S_STATUS_IN_DEV = SolutionStatus("1")

// S_STATUS_IN_DEV Assembled
const S_STATUS_ASSEMBLED = SolutionStatus("2")

// S_STATUS_IN_DEPLOYEMENT In Deployment
const S_STATUS_IN_DEPLOYEMENT = SolutionStatus("3")

// S_STATUS_DEPLOYED Deployed
const S_STATUS_DEPLOYED = SolutionStatus("4")

// S_PHASE_ACTIVATION
const S_PHASE_ACTIVATION = SolutionPhase("ACT")

// S_PHASE_DEVELOPMENT
const S_PHASE_DEVELOPMENT = SolutionPhase("DEV")

// S_PHASE_IMPORT
const S_PHASE_IMPORT = SolutionPhase("IMP")

// S_PHASE_DATA_UPDATE
const S_PHASE_DATA_UPDATE = SolutionPhase("POS")

// S_PHASE_STATUS_SUCCESSFUL
const S_PHASE_STATUS_SUCCESSFUL = SolutionPhaseStatus("S")

// S_PAHSE_STATUS_RUNNING
const S_PAHSE_STATUS_RUNNING = SolutionPhaseStatus("R")

// SolutionHeader information
//
// includes some status/text information
type SolutionHeader struct {
	ChangeDateTime     time.Time
	SolutionID         string
	SolutionName       string
	OriginSolutionID   string
	OriginSolutionName string
	Version            int64
	Status             SolutionStatus
	StatusText         string
	Phase              SolutionPhase
	PhaseStatus        SolutionPhaseStatus
	CanActivation      bool
	CanAssemble        bool
	CanDownload        bool
	NeedCreatePatch    bool
	IsSplitEnabled     bool
	IsRunningJob       bool
	IsCreatingPatch    bool
	// Help Text on PDI UI
	HelpText string
	// Is Solution Enabled
	Enabled bool
}

// IsRunningUploading file upload process
func (h *SolutionHeader) IsRunningUploading() bool {
	if h.Phase == S_PHASE_IMPORT && h.PhaseStatus == S_PAHSE_STATUS_RUNNING {
		return true
	}
	return false
}

// IsUploadingSuccessful file upload finished
func (h *SolutionHeader) IsUploadingSuccessful() bool {
	if h.Phase == S_PHASE_IMPORT && h.PhaseStatus == S_PHASE_STATUS_SUCCESSFUL {
		return true
	}
	return false
}

// IsRunningActivation process
//
// just used for deployment scenario
func (h *SolutionHeader) IsRunningActivation() bool {
	if h.Phase == S_PHASE_ACTIVATION && h.PhaseStatus == S_PAHSE_STATUS_RUNNING {
		return true
	}
	return false
}

// IsRunningDateUpdate process
func (h *SolutionHeader) IsRunningDataUpdate() bool {
	if h.Phase == S_PHASE_DATA_UPDATE && h.PhaseStatus == S_PAHSE_STATUS_RUNNING {
		return true
	}
	return false
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
	status := SolutionStatus(solutionHeader.Get("VERSION_STATUS").String())
	statusText := solutionHeader.Get("VERSION_STATUS_TEXT").String()
	helpText := solutionHeader.Get("HELP_TEXT").String()
	phase := SolutionPhase(solutionHeader.Get("PHASE").String())
	phaseStatus := SolutionPhaseStatus(solutionHeader.Get("PHASE_STATUS").String())
	enabled := solutionHeader.Get("IS_ENABLED").String() == "X"
	canActivation := solutionHeader.Get("EV_ACT_STATUS").String() == "X"
	canAssemble := solutionHeader.Get("EV_ASSEMBLE_STATUS").String() == "X"
	canDownload := solutionHeader.Get("EV_DOWNLOAD_STATUS").String() == "X"
	isRunningJob := solutionHeader.Get("EV_IS_SPLIT_JOB_RUNNING").String() == "X"
	isCreatingPatch := solutionHeader.Get("EV_IS_PATCH_JOB_RUNNING").String() == "X"
	isSplitEnabled := solutionHeader.Get("EV_IS_SPLIT_ENABLED").String() == "X"
	needCreatePatch := solutionHeader.Get("IS_PATCHSOL_REQUIRED").String() == "X"
	originSolutionID := solutionHeader.Get("ORIGIN_PROJECT_NAME").String()
	originSolutionName := solutionHeader.Get("ORIGIN_PROJECT_DESCRIPTION").String()

	return &SolutionHeader{
		ChangeDateTime:     changeDateTime,
		SolutionID:         solutionID,
		SolutionName:       solutionName,
		Version:            version,
		Status:             status,
		StatusText:         statusText,
		Phase:              phase,
		PhaseStatus:        phaseStatus,
		HelpText:           helpText,
		Enabled:            enabled,
		CanActivation:      canActivation,
		CanAssemble:        canAssemble,
		CanDownload:        canDownload,
		IsRunningJob:       isRunningJob,
		NeedCreatePatch:    needCreatePatch,
		IsCreatingPatch:    isCreatingPatch,
		IsSplitEnabled:     isSplitEnabled,
		OriginSolutionID:   originSolutionID,
		OriginSolutionName: originSolutionName,
	}

}

// IsHotFixNow func
//
// check solution if enable hot fix now
func (c *PDIClient) IsHotFixNow(solution string) (bool, error) {
	rt := false
	body, err := c.xrepRequestE("00163E0115B01DDFB194E54BB7202C9B", JSONObject{
		"IMPORTING": JSONObject{"IV_PRODUCT_NAME": solution},
	})

	if err != nil {
		return rt, err
	}

	rt = gjson.Get(body, "EXPORTING.ES_PRODUCT.VERSIONS.0.PV_OVERALL_STATUS").String() == "Deployed-In Correction"

	return rt, nil
}
