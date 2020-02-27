package pdiutil

import (
	"log"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type SolutionStatus string
type SolutionPhase string
type SolutionPhaseStatus string

// S_STATUS_IN_DEV Solution In Development
const S_STATUS_IN_DEV = SolutionStatus("1")

// S_STATUS_ASSEMBLED Assembled
const S_STATUS_ASSEMBLED = SolutionStatus("2")

// S_STATUS_IN_DEPLOYMENT In Deployment
const S_STATUS_IN_DEPLOYMENT = SolutionStatus("3")

// S_STATUS_DEPLOYED Deployed
const S_STATUS_DEPLOYED = SolutionStatus("4")

// S_PHASE_ACTIVATION
const S_PHASE_ACTIVATION = SolutionPhase("ACT")

// S_PHASE_ASSEMBLE
const S_PHASE_ASSEMBLE = SolutionPhase("EXP")

// S_PHASE_PATCH_CREATION
const S_PHASE_PATCH_CREATION = SolutionPhase("PTCH_CRT")

// S_PHASE_DEVELOPMENT
const S_PHASE_DEVELOPMENT = SolutionPhase("DEV")

// S_PHASE_IMPORT
const S_PHASE_IMPORT = SolutionPhase("IMP")

// S_PHASE_DATA_UPDATE
const S_PHASE_DATA_UPDATE = SolutionPhase("POS")

// S_PHASE_STATUS_SUCCESSFUL
const S_PHASE_STATUS_SUCCESSFUL = SolutionPhaseStatus("S")

// S_PHASE_STATUS_RUNNING
const S_PHASE_STATUS_RUNNING = SolutionPhaseStatus("R")

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
	if h.Phase == S_PHASE_IMPORT && h.PhaseStatus == S_PHASE_STATUS_RUNNING {
		return true
	}
	return false
}

// GetSolutionLatestAssembledVersion id
func (h *SolutionHeader) GetSolutionLatestAssembledVersion() int64 {
	if h.Status == S_STATUS_ASSEMBLED || h.IsCreatingPatch {
		return h.Version
	}
	return (h.Version - 1)
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
	if h.Phase == S_PHASE_ACTIVATION && h.PhaseStatus == S_PHASE_STATUS_RUNNING {
		return true
	}
	return false
}

// IsRunningDataUpdate process
func (h *SolutionHeader) IsRunningDataUpdate() bool {
	if h.Phase == S_PHASE_DATA_UPDATE && h.PhaseStatus == S_PHASE_STATUS_RUNNING {
		return true
	}
	return false
}

// PDILog Type
type PDILog struct {
	Sequence int
	Text     string
	Level    int
	Severity string
}

// GetSolutionLogs sync
func (c *PDIClient) GetSolutionLogs(solution string, version int64) []*PDILog {
	rt := []*PDILog{}
	solutionID := c.GetSolutionIDByString(solution)
	endpoint := "00163E0123D21EE092C936CC65A49BA4"
	payload := JSONObject{
		"IMPORTING": JSONObject{
			"IV_LANGUAGE":      "E",
			"IV_LOGLEVEL":      "4",
			"IV_SOLUTION_NAME": solutionID,
		},
	}
	respBody, err := c.xrepRequestE(endpoint, payload)

	if err != nil {
		panic(err)
	}

	versionList := gjson.Get(respBody, "EXPORTING.ES_SOLUTION_HEADER.SOLUTION_VERSION_LIST")

	for _, item := range versionList.Array() {
		versionID, _ := strconv.ParseInt(item.Get("VERSION_ID").String(), 10, 32)
		if version == versionID {
			for _, serverLog := range item.Get("SOLUTION_VERS_STATUS_LIST.0.SOLUTION_VERS_STATUS_LOG_LIST").Array() {
				seqNumber, _ := strconv.ParseInt(serverLog.Get("SEQUENCE_NUMBER").String(), 10, 32)
				logLevel, _ := strconv.ParseInt(serverLog.Get("LOGLEVEL").String(), 10, 32)
				severity := serverLog.Get("SEVERITY").String()
				text := serverLog.Get("TEXT").String()

				rt = append(rt, &PDILog{
					Sequence: int(seqNumber),
					Level:    int(logLevel),
					Text:     text,
					Severity: severity,
				})
			}
			break
		}
	}

	return rt
}

// GetSolutionStatus exported
func (c *PDIClient) GetSolutionStatus(solution string) *SolutionHeader {

	solutionID := c.GetSolutionIDByString(solution)
	endpoint := "00163E0123D21EE092C936CC65A49BA4"
	payload := JSONObject{
		"IMPORTING": JSONObject{
			"IV_LANGUAGE":      "E",
			"IV_LOGLEVEL":      "0",
			"IV_SOLUTION_NAME": solutionID,
		},
	}

	respBody, err := c.xrepRequestE(endpoint, payload)

	if err != nil {
		log.Println(err)
		log.Println("Some errors occurred when retrieve solution status information, will retry one time")
		// wait some seconds
		time.Sleep(DefaultPackageCheckInterval * time.Second)
		// retry once, maybe system is in busy
		respBody, err = c.xrepRequestE(endpoint, payload)
	}

	if err != nil {
		panic(err)
	}

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
	// fix for new c4c release
	canDownload := ((phase == S_PHASE_ASSEMBLE) && (phaseStatus == S_PHASE_STATUS_SUCCESSFUL))
	isRunningJob := solutionHeader.Get("EV_IS_SPLIT_JOB_RUNNING").String() == "X"
	isCreatingPatch := (solutionHeader.Get("EV_IS_PATCH_JOB_RUNNING").String() == "X" || phase == S_PHASE_PATCH_CREATION)
	isSplitEnabled := solutionHeader.Get("EV_IS_SPLIT_ENABLED").String() == "X"
	needCreatePatch := status == S_STATUS_ASSEMBLED
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
