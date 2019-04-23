package pdiutil

import (
	"github.com/tidwall/gjson"
)

// BPUserNameMapping func
//
// name is tech id, value is business user id
// for cache usage
type BPUserNameMapping = map[string]string

var cachedMapping = BPUserNameMapping{}

// GetAUserIDNameByTechID func
func (c *PDIClient) GetAUserIDNameByTechID(techID string) string {
	rt := ""
	if techID != "" {

		if cached, hit := cachedMapping[techID]; hit {
			rt = cached
		} else {
			rt = c.GetBPUserNameByTechID(append([]string{}, techID))[techID]
			// >> cache value
			if rt != "" {
				cachedMapping[techID] = rt
			} else {
				// failback, avoid query again
				cachedMapping[techID] = "Unknown User TechID: " + techID
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
