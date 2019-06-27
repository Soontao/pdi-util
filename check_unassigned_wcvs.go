package pdiutil

import (
	"encoding/xml"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

const (
	// DefaultDownloadConcurrent value
	//
	// used to limit the download concurrent
	DefaultDownloadConcurrent = 20
	// SuffixWCV
	//
	// Work Center View suffix
	SuffixWCV = "WCVIEW.uiwocview"
	// SuffixWoC
	//
	// Work Center suffix
	SuffixWoC = "WCF.uiwoc"
)

var extractWCVReferenceReg = regexp.MustCompile(`<uxc:EmbeddedComponent.*?targetComponentID="(.*?)".*?/>
`)

// WCVAssignCheckResult type
type WCVAssignCheckResult struct {
	WCCount            int
	WCVCount           int
	AssignedWCVCount   int
	UnAssignedWCVCount int
	UnAssignedWCVs     []string
}

// FindUnAssignedWCV api
func (c *PDIClient) FindUnAssignedWCV(solution string) *WCVAssignCheckResult {
	files := c.GetSolutionXrepFileList(solution)
	allWorkCenters := []string{}
	allWorkCenterViews := []string{}
	usedWorkCenterViews := mapset.NewSet()
	unusedWorkCenterViews := []string{}
	for _, f := range files {
		if strings.HasSuffix(f, SuffixWoC) {
			allWorkCenters = append(allWorkCenters, f)
		} else if strings.HasSuffix(f, SuffixWCV) {
			allWorkCenterViews = append(allWorkCenterViews, f)
		}
	}

	wcSources := c.fetchSources(allWorkCenters, DefaultDownloadConcurrent)

	// extract used WCV in WoC
	for _, wcSource := range wcSources {
		wcStruct := &UXComponent{}
		// ignore parse error
		xml.Unmarshal(wcSource.Source, wcStruct)

		for _, ec := range wcStruct.EmbeddedComponents.EmbeddedComponent {
			if strings.HasSuffix(ec.TargetComponentID, SuffixWCV) {
				usedWorkCenterViews.Add(ec.TargetComponentID)
			}
		}

	}

	for _, wc := range allWorkCenterViews {
		if !usedWorkCenterViews.Contains(wc) {
			unusedWorkCenterViews = append(unusedWorkCenterViews, wc)
		}
	}

	return &WCVAssignCheckResult{
		WCCount:            len(allWorkCenters),
		WCVCount:           len(allWorkCenterViews),
		AssignedWCVCount:   len(usedWorkCenterViews.ToSlice()),
		UnAssignedWCVCount: len(unusedWorkCenterViews),
		UnAssignedWCVs:     unusedWorkCenterViews,
	}

}

// UXComponent for WoC
type UXComponent struct {
	XMLName                         xml.Name `xml:"UXComponent"`
	Text                            string   `xml:",chardata"`
	Base                            string   `xml:"base,attr"`
	Uxv                             string   `xml:"uxv,attr"`
	Uxc                             string   `xml:"uxc,attr"`
	ID                              string   `xml:"id,attr"`
	Name                            string   `xml:"name,attr"`
	ComponentType                   string   `xml:"componentType,attr"`
	AuthorizationClassificationCode string   `xml:"AuthorizationClassificationCode,attr"`
	HelpId                          string   `xml:"helpId,attr"`
	EnableBackendOperationsClubbing string   `xml:"enableBackendOperationsClubbing,attr"`
	DesigntimeVersion               string   `xml:"designtimeVersion,attr"`
	ModelVersion                    string   `xml:"modelVersion,attr"`
	UseUIController                 string   `xml:"useUIController,attr"`
	Xmlns                           string   `xml:"xmlns,attr"`
	Interface                       struct {
		Text          string `xml:",chardata"`
		ID            string `xml:"id,attr"`
		Configuration struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"Configuration"`
	} `xml:"Interface"`
	Implementation struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"id,attr"`
		DataModel struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"id,attr"`
			Name        string `xml:"name,attr"`
			PropertyBag struct {
				Text     string `xml:",chardata"`
				ID       string `xml:"id,attr"`
				Property []struct {
					Text  string `xml:",chardata"`
					ID    string `xml:"id,attr"`
					Name  string `xml:"name,attr"`
					Value string `xml:"value,attr"`
				} `xml:"Property"`
			} `xml:"PropertyBag"`
			DataField []struct {
				Text        string `xml:",chardata"`
				ID          string `xml:"id,attr"`
				Name        string `xml:"name,attr"`
				Type        string `xml:"type,attr"`
				PropertyBag struct {
					Text     string `xml:",chardata"`
					ID       string `xml:"id,attr"`
					Property []struct {
						Text  string `xml:",chardata"`
						ID    string `xml:"id,attr"`
						Name  string `xml:"name,attr"`
						Value string `xml:"value,attr"`
					} `xml:"Property"`
				} `xml:"PropertyBag"`
			} `xml:"DataField"`
			Structure []struct {
				Text        string `xml:",chardata"`
				ID          string `xml:"id,attr"`
				Name        string `xml:"name,attr"`
				PropertyBag struct {
					Text     string `xml:",chardata"`
					ID       string `xml:"id,attr"`
					Property []struct {
						Text  string `xml:",chardata"`
						ID    string `xml:"id,attr"`
						Name  string `xml:"name,attr"`
						Value string `xml:"value,attr"`
					} `xml:"Property"`
				} `xml:"PropertyBag"`
				DataField []struct {
					Text         string `xml:",chardata"`
					ID           string `xml:"id,attr"`
					Name         string `xml:"name,attr"`
					Type         string `xml:"type,attr"`
					InitialValue string `xml:"initialValue,attr"`
					PropertyBag  struct {
						Text     string `xml:",chardata"`
						ID       string `xml:"id,attr"`
						Property []struct {
							Text  string `xml:",chardata"`
							ID    string `xml:"id,attr"`
							Name  string `xml:"name,attr"`
							Value string `xml:"value,attr"`
						} `xml:"Property"`
					} `xml:"PropertyBag"`
				} `xml:"DataField"`
			} `xml:"Structure"`
		} `xml:"DataModel"`
		ModelChanges string `xml:"ModelChanges"`
	} `xml:"Implementation"`
	Navigation struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"Navigation"`
	CenterStructure struct {
		Text  string `xml:",chardata"`
		ID    string `xml:"id,attr"`
		Title struct {
			Text       string `xml:",chardata"`
			ID         string `xml:"id,attr"`
			TextPoolId string `xml:"textPoolId,attr"`
		} `xml:"Title"`
		ViewSwitches struct {
			Text       string `xml:",chardata"`
			ID         string `xml:"id,attr"`
			ViewSwitch []struct {
				Chardata    string `xml:",chardata"`
				ID          string `xml:"id,attr"`
				PropertyBag struct {
					Text     string `xml:",chardata"`
					ID       string `xml:"id,attr"`
					Property struct {
						Text  string `xml:",chardata"`
						ID    string `xml:"id,attr"`
						Name  string `xml:"name,attr"`
						Value string `xml:"value,attr"`
					} `xml:"Property"`
				} `xml:"PropertyBag"`
				Text struct {
					Text       string `xml:",chardata"`
					ID         string `xml:"id,attr"`
					TextPoolId string `xml:"textPoolId,attr"`
				} `xml:"Text"`
				SubViewSwitches struct {
					Text       string `xml:",chardata"`
					ID         string `xml:"id,attr"`
					ViewSwitch []struct {
						Chardata    string `xml:",chardata"`
						ID          string `xml:"id,attr"`
						EmbedName   string `xml:"embedName,attr"`
						PropertyBag struct {
							Text     string `xml:",chardata"`
							ID       string `xml:"id,attr"`
							Property struct {
								Text  string `xml:",chardata"`
								ID    string `xml:"id,attr"`
								Name  string `xml:"name,attr"`
								Value string `xml:"value,attr"`
							} `xml:"Property"`
						} `xml:"PropertyBag"`
						Text struct {
							Text       string `xml:",chardata"`
							ID         string `xml:"id,attr"`
							TextPoolId string `xml:"textPoolId,attr"`
						} `xml:"Text"`
					} `xml:"ViewSwitch"`
				} `xml:"SubViewSwitches"`
			} `xml:"ViewSwitch"`
		} `xml:"ViewSwitches"`
		Attributes struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id,attr"`
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"Attributes"`
	} `xml:"CenterStructure"`
	EmbeddedComponents struct {
		Text              string `xml:",chardata"`
		ID                string `xml:"id,attr"`
		EmbeddedComponent []struct {
			Text              string `xml:",chardata"`
			ID                string `xml:"id,attr"`
			EmbedName         string `xml:"embedName,attr"`
			TargetComponentID string `xml:"targetComponentID,attr"`
			LazyLoad          string `xml:"lazyLoad,attr"`
		} `xml:"EmbeddedComponent"`
	} `xml:"EmbeddedComponents"`
	Description struct {
		Text       string `xml:",chardata"`
		ID         string `xml:"id,attr"`
		TextPoolId string `xml:"textPoolId,attr"`
	} `xml:"Description"`
	Prerequisites      string `xml:"Prerequisites"`
	TechnicalConflicts string `xml:"TechnicalConflicts"`
	SoDConflicts       string `xml:"SoDConflicts"`
	TextPool           struct {
		Text      string `xml:",chardata"`
		TextBlock struct {
			Text            string `xml:",chardata"`
			Language        string `xml:"language,attr"`
			MasterLanguage  string `xml:"masterLanguage,attr"`
			CurrentLanguage string `xml:"currentLanguage,attr"`
			TextPoolEntry   []struct {
				Text         string `xml:",chardata"`
				ID           string `xml:"id,attr"`
				TextUuid     string `xml:"textUuid,attr"`
				AttrText     string `xml:"text,attr"`
				TextCategory string `xml:"textCategory,attr"`
			} `xml:"TextPoolEntry"`
		} `xml:"TextBlock"`
	} `xml:"TextPool"`
}
