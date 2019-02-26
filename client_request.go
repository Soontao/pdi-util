package main

import (
	"fmt"
	"strings"

	"github.com/imroc/req"
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
		panic(fmt.Errorf("error with %s, with payload %+v", err, payload))
	}
	respBody, _ := resp.ToString()
	return respBody
}
