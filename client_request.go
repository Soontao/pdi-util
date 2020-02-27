package pdiutil

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

func (c *PDIClient) path(path string) string {
	if !strings.HasPrefix(path, "/") {
		panic("path must start with /")
	}
	if strings.HasPrefix(c.hostname, "https://") {
		c.hostname = strings.TrimPrefix(c.hostname, "https://")
	}
	return fmt.Sprintf("https://%s%s", strings.TrimSpace(c.hostname), path)
}

func (c *PDIClient) xrepPath() string {
	return c.path("/sap/ap/xrep/json3")
}

func (c *PDIClient) query(fm string) req.QueryParam {
	return req.QueryParam{"stateful": "0", "sap-client": c.sapClient, "fm": fm}
}

func (c *PDIClient) xrepRequest(code string, payload interface{}) string {
	url := c.xrepPath()
	query := c.query(code)
	resp, err := req.Post(url, req.BodyJSON(payload), query)
	if err != nil {
		panic(fmt.Errorf("Request endpoint: %v with error: %s", code, err))
	}
	respBody, _ := resp.ToString()
	return respBody
}

func (c *PDIClient) xrepRequestE(code string, payload interface{}) (res string, err error) {

	checkInterval := time.Second * DefaultPackageCheckInterval

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	res = c.xrepRequest(code, payload)

	for { // retry if server on maintenance
		if strings.Contains(res, "System Maintenance / Temporary System Downtime") {
			log.Println("Server is under maintenance, please wait resume")
			time.Sleep(checkInterval)
			res = c.xrepRequest(code, payload)
		} else {
			break
		}
	}

	success := gjson.Get(res, "EXPORTING.EV_SUCCESS").String() == "X"

	if !success {
		err = fmt.Errorf(gjson.Get(res, "EXPORTING.ET_MESSAGES").String())
	}

	return res, err
}
