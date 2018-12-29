package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// XrepFileAttrs
type XrepFileAttrs struct {
	FileName      string
	FilePath      string
	EntityType    string
	CreatedBy     string
	LastChangedBy string
	CreatedOn     string
	LastChangedOn string
	FileSize      string
	MimeType      string
	IsFolder      bool
	Solution      string
	Branch        string
	ActiveFlag    string
	Active        bool
	ContID        string
	IsLink        bool
}

// GetSolutionFileAttrs from xrep
func (c *PDIClient) GetSolutionFileAttrs(solutionName string) map[string]XrepFileAttrs {
	rt := map[string]XrepFileAttrs{}
	payload := map[string]interface{}{
		"IMPORTING": map[string]interface{}{
			"IT_FILTER":               []string{},
			"IV_LAST_SHIPPED_VERSION": "",
			"IV_PATH":                 fmt.Sprintf("/%sMAIN", solutionName),
			"IV_RECURSIVELY":          "X",
			"IV_VIRTUAL_VIEW":         "X",
			"IV_WITH_ATTRIBUTE":       "X",
		},
	}

	respBody := c.xrepRequest("00163E0115B01DDFB194E54BB722CC9B", payload)
	fileList := gjson.Get(respBody, "EXPORTING.ET_CONT_INFO").Array()
	for _, fileAttrJSON := range fileList {
		xrepFileAttr := XrepFileAttrs{}
		xrepFileAttr.FileName = fileAttrJSON.Get("FILE_PATH").String()
		xrepFileAttr.FilePath = fileAttrJSON.Get("FULL_PATH").String()
		for _, attr := range fileAttrJSON.Get("FILE_ATTRS").Array() {
			attrName := attr.Get("NAME").String()
			attrValue := attr.Get("VALUE").String()

			switch attrName {
			case "~ENTITY_TYPE":
				xrepFileAttr.EntityType = attrValue
			case "~CREATED_BY":
				xrepFileAttr.CreatedBy = attrValue
			case "~LAST_CHANGED_BY":
				xrepFileAttr.LastChangedBy = attrValue
			case "~LAST_CHANGED_ON":
				xrepFileAttr.LastChangedOn = attrValue
			case "~CREATED_ON":
				xrepFileAttr.CreatedOn = attrValue
			case "~FILE_SIZE":
				xrepFileAttr.FileSize = attrValue
			case "~MIME_TYPE":
				xrepFileAttr.MimeType = attrValue
			case "~IS_FOLDER":
				xrepFileAttr.IsFolder = (attrValue == "X")
			case "~SOLUTION":
				xrepFileAttr.Solution = attrValue
			case "~BRANCH":
				xrepFileAttr.Branch = attrValue
			case "~ACTIVE_FLAG":
				// ActiveFlag maybe 'A' or 'I'. 'I' means inactive
				xrepFileAttr.Active = (attrValue == "A")
				xrepFileAttr.ActiveFlag = attrValue
			case "~CONT_ID":
				xrepFileAttr.ContID = attrValue
			case "~IS_LINK":
				xrepFileAttr.IsLink = (attrValue == "X")

			}
		}
		rt[xrepFileAttr.FileName] = xrepFileAttr
	}

	return rt
}
