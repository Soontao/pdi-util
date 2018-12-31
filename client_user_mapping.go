package main

import (
	"github.com/tidwall/gjson"
)

// BPUserNameMapping
// name is tech id, value is business user id
type BPUserNameMapping = map[string]string

var cachedMapping = BPUserNameMapping{}

// GetAUserIDNameByTechID
func (c *PDIClient) GetAUserIDNameByTechID(techID string) string {
	rt := ""
	if techID != "" {
		cached := cachedMapping[techID]
		if cached != "" && cached != "null" {
			rt = cached
		} else {
			rt = c.GetBPUserNameByTechID(append([]string{}, techID))[techID]
			// >> cache value
			if rt != "" {
				cachedMapping[techID] = rt
			} else {
				cachedMapping[techID] = "null"
			}
		}
	}
	return rt
}

// GetBPUserNameByTechID
func (c *PDIClient) GetBPUserNameByTechID(techIDs []string) BPUserNameMapping {
	rt := BPUserNameMapping{}
	if len(techIDs) > 0 {
		names := []JSONObject{}
		for _, techID := range techIDs {
			names = append(names, JSONObject{"BAPIBNAME": techID})
		}
		payload := JSONObject{
			"IMPORTING": JSONObject{
				"IT_BAPIBNAME": names,
			},
		}
		resp := c.xrepRequest("0000000000011EE19DFE474CB9888B14", payload)
		mappingInformation := gjson.Get(resp, "EXPORTING.ET_BAPIALIAS").Array()
		for _, mapping := range mappingInformation {
			rt[mapping.Get("BAPIBNAME").String()] = mapping.Get("BAPIALIAS").String()
		}
	}
	return rt
}
