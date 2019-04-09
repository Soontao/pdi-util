package pdiutil

import "encoding/xml"

type JSONObject = map[string]interface{}

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

// Project struct for pdi project xml
type Project struct {
	XMLName        xml.Name `xml:"Project"`
	Text           string   `xml:",chardata"`
	DefaultTargets string   `xml:"DefaultTargets,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	PropertyGroup  []struct {
		Text          string `xml:",chardata"`
		Condition     string `xml:"Condition,attr"`
		SchemaVersion string `xml:"SchemaVersion"`
		ProjectGUID   string `xml:"ProjectGuid"`
		ProjectType   string `xml:"ProjectType"`
		Configuration struct {
			Text      string `xml:",chardata"`
			Condition string `xml:"Condition,attr"`
		} `xml:"Configuration"`
		Name                      string `xml:"Name"`
		RootNamespace             string `xml:"RootNamespace"`
		RepositoryNamespace       string `xml:"RepositoryNamespace"`
		RuntimeNamespacePrefix    string `xml:"RuntimeNamespacePrefix"`
		RepositoryRootFolder      string `xml:"RepositoryRootFolder"`
		DefaultProcessComponent   string `xml:"DefaultProcessComponent"`
		DevelopmentPackage        string `xml:"DevelopmentPackage"`
		XRepSolution              string `xml:"XRepSolution"`
		BCSourceFolderInXRep      string `xml:"BCSourceFolderInXRep"`
		ProjectSourceFolderinXRep string `xml:"ProjectSourceFolderinXRep"`
		DeploymentUnit            string `xml:"DeploymentUnit"`
		CompilerVersion           string `xml:"CompilerVersion"`
		OutputPath                string `xml:"OutputPath"`
	} `xml:"PropertyGroup"`
	ItemGroup []struct {
		Text   string `xml:",chardata"`
		Folder []struct {
			Text    string `xml:",chardata"`
			Include string `xml:"Include,attr"`
		} `xml:"Folder"`
		BCSet []struct {
			Text    string `xml:",chardata"`
			Include string `xml:"Include,attr"`
			SubType string `xml:"SubType"`
		} `xml:"BCSet"`
		Content []struct {
			Text          string `xml:",chardata"`
			Include       string `xml:"Include,attr"`
			DependentUpon string `xml:"DependentUpon"`
			SubType       string `xml:"SubType"`
			IsHidden      string `xml:"IsHidden"`
		} `xml:"Content"`
	} `xml:"ItemGroup"`
	Import struct {
		Text    string `xml:",chardata"`
		Project string `xml:"Project,attr"`
	} `xml:"Import"`
}
