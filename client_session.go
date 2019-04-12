package pdiutil

import (
	"fmt"
	"log"

	"github.com/tidwall/gjson"
)

const DefaultRelease = "1902"

// GetSessionID from server
// will be triggerred at login
func (c *PDIClient) GetSessionID(release string) (sessionID string, err error) {

	if release == "" {
		release = DefaultRelease
	}

	payload := JSONObject{
		"IMPORTING": JSONObject{
			"IV_TIMEOUT": 31536000,
			"IS_LOGON_COMPONENT": JSONObject{
				"BASE_COMPONENT": "SAP_LEAP",
				"BYD_RELEASE":    release,
			},
		},
	}

	resp := c.xrepRequest("00163E0115B01DDFB194E54BB7206C9B", payload)
	success := gjson.Get(resp, "EXPORTING.EV_SUCCESS").String() == "X"

	if success {
		sessionID = gjson.Get(resp, "EXPORTING.EV_SID").String()
		c.sessionID = sessionID
		log.Printf("Retrive session id: %v", sessionID)
	} else {
		s := ""
		for _, e := range gjson.Get(resp, "EXPORTING.ET_MESSAGES").Array() {
			s += e.Get("TEXT").String()
			s += ", "
		}
		err = fmt.Errorf(s)
	}

	return sessionID, err
}
