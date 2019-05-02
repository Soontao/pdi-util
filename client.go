package pdiutil

import (
	"fmt"
	"log"
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
	// tenant release version
	release string
	// sap client id
	sapClient string
	// tech user id
	ivUser string
	// exit Code
	exitCode int
	// user session id
	sessionID string
}

// GetExitCode for client
func (c *PDIClient) GetExitCode() int {
	return c.exitCode
}

func (c *PDIClient) login() (*PDIClient, error) {
	url := c.path("/sap/ap/ui/login")
	// > fetch cookie & client infomartions
	query := req.QueryParam{"saml2": "disabled"}
	resp, err := req.Get(url, query)
	if err != nil {
		return nil, err
	}
	respBody := &LoginResponse{}
	err = resp.ToXML(respBody)

	if err != nil {
		return nil, err
	}

	// > login
	param := req.Param{}
	for _, aData := range respBody.Data.Element {
		param[aData.Name] = aData.Value
	}
	param["sap-alias"] = strings.TrimSpace(c.username)
	param["sap-system-login-oninputprocessing"] = "onLogin"
	param["sap-password"] = strings.TrimSpace(c.password)
	delete(param, "sap-user")
	c.sapClient = param["sap-client"].(string)

	resp, err = req.Post(url, param)
	if err != nil {
		return nil, err
	}

	if resp.Response().Header.Get("content-type") == "text/xml; charset=utf-8" {
		resp.ToXML(respBody)
		isError := respBody.Messages.Message[0].Type == "error"
		msg := respBody.Messages.Message[0].AttrText
		session := respBody.Sessions.Session.Terminal
		if isError {
			return nil, fmt.Errorf("%s%s", msg, session)
		}
	}

	if _, err := c.GetSessionID(c.release); err != nil {
		return nil, err
	}

	c.getIvUser()

	return c, nil
}

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
	log.Printf("Login to %v as %v(%v)", c.hostname, c.ivUser, c.username)
	return c
}

func ensure(v interface{}, name string) {
	if v == "" {
		panic(fmt.Errorf("You must give out the %s", name))
	}
}

// NewPDIClient instance
func NewPDIClient(username, password, hostname, release string) (c *PDIClient, err error) {

	ensure(username, "username")
	ensure(password, "password")
	ensure(hostname, "hostname")
	ensure(release, "release")
	c = &PDIClient{username: username, password: password, hostname: hostname, release: release, exitCode: 0}

	_, err = c.login()

	return c, err
}
