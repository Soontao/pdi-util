package pdiutil

import (
	"encoding/base64"
	"encoding/xml"

	"github.com/tidwall/gjson"
)

// QueryBOInformation func
func (c *PDIClient) QueryBOInformation(solution string) *BOInformation {
	namespace := c.GetSolutionNamespace(solution)
	rt := &BOInformation{}
	payload := JSONObject{"IMPORTING": JSONObject{
		"IV_WITH_PSM_FILTER": "",
		"IS_QUERY_OPTION":    nil,
		"IT_SELECTION_PARAMETER": []JSONObject{
			{
				"ATTRIBUTE_NAME": "NAME_KEY-NAMESPACE",
				"HIGH":           nil,
				"LOW":            namespace,
				"OPTION":         "EQ",
				"SIGN":           "I",
			},
			{
				"ATTRIBUTE_NAME": "KEY-VERSION_ID",
				"HIGH":           nil,
				"LOW":            "0",
				"OPTION":         "EQ",
				"SIGN":           "I",
			},

			{
				"ATTRIBUTE_NAME": "KEY-VERSION_ID",
				"HIGH":           nil,
				"LOW":            "1",
				"OPTION":         "EQ",
				"SIGN":           "I",
			},
		},
	}}
	res := c.xrepRequest("00163E0115B01DDFB194E54BB721EC9B", payload)

	base64ResultSet := gjson.Get(res, "EXPORTING.EV_RESULT_SET").String()

	bytes, err := base64.StdEncoding.DecodeString(base64ResultSet)

	if err != nil {
		panic(err)
	}

	if err := xml.Unmarshal(bytes, rt); err != nil {
		panic(err)
	}

	return rt

}

// QueryBCInformation
func (c *PDIClient) QueryBCInformation(solution string) *BCInformation {
	namespace := c.GetSolutionNamespace(solution)
	rt := &BCInformation{}
	payload := JSONObject{"IMPORTING": JSONObject{
		"IV_WITH_PSM_FILTER": "",
		"IS_QUERY_OPTION":    nil,
		"IT_SELECTION_PARAMETER": []JSONObject{
			{
				"ATTRIBUTE_NAME": "NAME_KEY-NAMESPACE",
				"HIGH":           nil,
				"LOW":            namespace,
				"OPTION":         "EQ",
				"SIGN":           "I",
			},
			{
				"ATTRIBUTE_NAME": "KEY-VERSION_ID",
				"HIGH":           nil,
				"LOW":            "0",
				"OPTION":         "EQ",
				"SIGN":           "I",
			},

			{
				"ATTRIBUTE_NAME": "KEY-VERSION_ID",
				"HIGH":           nil,
				"LOW":            "1",
				"OPTION":         "EQ",
				"SIGN":           "I",
			},
		},
	}}

	res := c.xrepRequest("00163E0115B01DDFB194E54BB721EC9B", payload)

	base64ResultSet := gjson.Get(res, "EXPORTING.EV_RESULT_SET").String()

	bytes, err := base64.StdEncoding.DecodeString(base64ResultSet)

	if err != nil {
		panic(err)
	}

	if err := xml.Unmarshal(bytes, rt); err != nil {
		panic(err)
	}

	return rt

}
