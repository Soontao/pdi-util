package pdiutil

import (
	"log"
	"path/filepath"
	"time"

	"github.com/tidwall/gjson"
)

type XrepActivateCheckResult struct {
	File                string
	FilePath            string
	LastChangedBy       string
	LastChangedByUserID string
	LastChangedOn       string
	LastChangedOnTime   time.Time
}

type XrepAddionalActiveCheckResult struct {
	FileName string
	FilePath string
	Active   bool
}

func (c *PDIClient) checkAddionalActiveStatus(solutionName string) map[string]XrepAddionalActiveCheckResult {
	rt := map[string]XrepAddionalActiveCheckResult{}
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IV_SOLUTION_PREFIX": solutionName,
		},
	}

	respBody := c.xrepRequest("0000000000011EE19CCE30BEB4D01B74", payload)
	contentsStatus := gjson.Get(respBody, "EXPORTING.ET_CONTENT_STATUS").Array()

	for _, contentStatus := range contentsStatus {
		filePath := contentStatus.Get("FILE_NAME").String()
		_, fileName := filepath.Split(filePath)
		runtimeStatus := contentStatus.Get("STATUS_RUN_TIME").String()
		active := false
		switch runtimeStatus {
		case "0":
			active = false
		case "1":
			active = false
		case "2":
			active = true
		default:
			log.Printf("For file %s, have unexpected Runtime Status %s", filePath, runtimeStatus)
		}
		rt[fileName] = XrepAddionalActiveCheckResult{fileName, filePath, active}
	}

	return rt
}

// CheckInActiveFilesAPI return all in active files
func (c *PDIClient) CheckInActiveFilesAPI(solutionName string) []XrepActivateCheckResult {
	rt := []XrepActivateCheckResult{}

	fileAttrs := c.GetSolutionFileAttrs(solutionName)

	addtionalActiveCheckResult := c.checkAddionalActiveStatus(solutionName)

	// merge addional active check result
	for _, addionalCheck := range addtionalActiveCheckResult {

		if fileAttr, existed := fileAttrs[addionalCheck.FileName]; existed {
			if fileAttr.Active && (!addionalCheck.Active) {
				fileAttr.Active = false
			}
		}

	}

	for _, fileAttr := range fileAttrs {
		if !fileAttr.Active {
			f := XrepActivateCheckResult{}
			f.File = fileAttr.FileName
			f.FilePath = fileAttr.FilePath
			f.LastChangedBy = fileAttr.LastChangedBy
			f.LastChangedByUserID = c.GetAUserIDNameByTechID(f.LastChangedBy)
			f.LastChangedOn = fileAttr.LastChangedOn
			f.LastChangedOnTime = ParseXrepDateString(f.LastChangedOn)
			rt = append(rt, f)
		}
	}

	return rt
}
