package client

import "encoding/xml"

// LoginResponse struct
type LoginResponse struct {
	XMLName  xml.Name `xml:"Login"`
	Text     string   `xml:",chardata"`
	Messages struct {
		Text    string `xml:",chardata"`
		Message []struct {
			Text     string `xml:",chardata"`
			Type     string `xml:"type,attr"`
			AttrText string `xml:"text,attr"`
		} `xml:"Message"`
	} `xml:"Messages"`
	Actions struct {
		Text   string `xml:",chardata"`
		Action struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"Action"`
	} `xml:"Actions"`
	SimpleTypes struct {
		Text       string `xml:",chardata"`
		SimpleType struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			Enum []struct {
				Text        string `xml:",chardata"`
				AttrText    string `xml:"text,attr"`
				Value       string `xml:"value,attr"`
				JaproSuffix string `xml:"japroSuffix,attr"`
			} `xml:"Enum"`
		} `xml:"SimpleType"`
	} `xml:"SimpleTypes"`
	Data struct {
		Text    string `xml:",chardata"`
		Element []struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"Element"`
	} `xml:"Data"`
	Sessions struct {
		Text    string `xml:",chardata"`
		Session struct {
			Text     string `xml:",chardata"`
			Client   string `xml:"client,attr"`
			User     string `xml:"user,attr"`
			Terminal string `xml:"terminal,attr"`
			Time     string `xml:"time,attr"`
		} `xml:"Session"`
	} `xml:"Sessions"`
	Config struct {
		Text      string `xml:",chardata"`
		Parameter []struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"Parameter"`
	} `xml:"Config"`
}

type IvUserRequest struct {
	Importing IvUserImporting `json:"IMPORTING"`
}

type IvUserImporting struct {
	IvAlias IvAlias `json:"IV_ALIAS"`
}

type IvAlias struct {
	Useralias string `json:"USERALIAS"`
}
