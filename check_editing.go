package main

import (
	"fmt"
	"path/filepath"

	"github.com/tidwall/gjson"
)

// EtLock information
type EtLock struct {
	FileName string
	FilePath string
	EditBy   string
	EditOn   string
}

// CheckLockedFilesAPI to get locked files
func (c *PDIClient) CheckLockedFilesAPI(solution string) []EtLock {
	rt := []EtLock{}
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IT_PATH": []string{fmt.Sprintf("%sBC/SRC", solution), fmt.Sprintf("%sMAIN/SRC", solution)},
		},
	}

	respBody := c.xrepRequest("00163E0115B01DDFB194E54BB722EC9B", payload)
	lockInformation := gjson.Get(respBody, "EXPORTING.ET_LOCK").Array()
	for _, lock := range lockInformation {
		l := EtLock{}
		l.FilePath = lock.Get("FILEPATH").String()
		_, l.FileName = filepath.Split(l.FilePath)
		l.EditBy = lock.Get("EDIT_BY").String()
		l.EditOn = lock.Get("EDIT_ON").String()
		rt = append(rt, l)
	}

	return rt
}
