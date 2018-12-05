package main

import (
	"fmt"
	"strings"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// PDIClient for programming
type PDIClient struct {
	// pdi user name
	username string
	// pdi user password
	password string
	// pdi tenant hostname
	hostname string
	// sap client id
	sapClient string
	// tech user id
	ivUser string
	// exit Code
	exitCode int
}

func (c *PDIClient) path(path string) string {
	if !strings.HasPrefix(path, "/") {
		panic("path must start with /")
	}
	return fmt.Sprintf("https://%s%s", c.hostname, path)
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
		panic(nil)
	}
	respBody, _ := resp.ToString()
	success := gjson.Get(respBody, "EXPORTING.EV_SUCCESS").String() == "X"
	if !success {
		message := gjson.Get(respBody, "EXPORTING.ET_MESSAGES").String()
		panic(message)
	}
	return respBody
}

func (c *PDIClient) login() *PDIClient {
	url := c.path("/sap/ap/ui/login")
	// > fetch cookie & client infomartions
	query := req.QueryParam{"saml2": "disabled"}
	resp, err := req.Get(url, query)
	if err != nil {
		panic(err)
	}
	respBody := &LoginResponse{}
	err = resp.ToXML(respBody)
	if err != nil {
		panic(err)
	}

	// > login
	param := req.Param{}
	for _, aData := range respBody.Data.Element {
		param[aData.Name] = aData.Value
	}
	param["sap-alias"] = c.username
	param["sap-system-login-oninputprocessing"] = "onLogin"
	param["sap-password"] = c.password
	delete(param, "sap-user")
	c.sapClient = param["sap-client"].(string)

	resp, err = req.Post(url, param)
	if err != nil {
		panic(err)
	}

	if resp.Response().Header.Get("content-type") == "text/xml; charset=utf-8" {
		resp.ToXML(respBody)
		isError := respBody.Messages.Message[0].Type == "error"
		msg := respBody.Messages.Message[0].AttrText
		session := respBody.Sessions.Session.Terminal
		if isError {
			panic(fmt.Sprintf("%s%s", msg, session))
		}
	}

	return c.getIvUser()
}

// Destroy PDI session
func (c *PDIClient) getIvUser() *PDIClient {
	url := c.path("/sap/ap/xrep/json_pdi")
	query := c.query("0000000000011ED19CEC5BA760AEE530")
	reqBody := IvUserRequest{IvUserImporting{IvAlias{strings.ToUpper(c.username)}}}
	resp, err := req.Post(url, req.BodyJSON(reqBody), query)
	if err != nil {
		panic(err)
	}
	respBody, _ := resp.ToString()
	c.ivUser = gjson.Get(respBody, "EXPORTING.EV_USER").String()
	return c
}

func ensure(v interface{}, name string) {
	if v == "" {
		panic(fmt.Sprintf("You must give out the %s!", name))
	}
}

// NewPDIClient instance
func NewPDIClient(username, password, hostname string) *PDIClient {
	ensure(username, "username")
	ensure(password, "password")
	ensure(hostname, "hostname")
	rt := &PDIClient{username: username, password: password, hostname: hostname, exitCode: 0}
	rt.login()
	return rt
}
