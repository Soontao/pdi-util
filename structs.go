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

// BOInformation struct
type BOInformation struct {
	XMLName xml.Name `xml:"abap"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Asx     string   `xml:"asx,attr"`
	Values  struct {
		Text      string `xml:",chardata"`
		RESULTSET struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text      string `xml:",chardata"`
				NODEID    string `xml:"NODE_ID"`
				NAME      string `xml:"NAME"`
				NAMESPACE string `xml:"NAMESPACE"`
				UUID      string `xml:"UUID"`
				PROXYNAME string `xml:"PROXYNAME"`
				NODE      struct {
					Text string `xml:",chardata"`
					Item struct {
						Text           string `xml:",chardata"`
						NODEID         string `xml:"NODE_ID"`
						NAME           string `xml:"NAME"`
						UUID           string `xml:"UUID"`
						PROXYNAME      string `xml:"PROXYNAME"`
						PARENTNODENAME string `xml:"PARENT_NODE_NAME"`
						WRITEACCESS    string `xml:"WRITE_ACCESS"`
						CATEGORY       string `xml:"CATEGORY"`
					} `xml:"item"`
				} `xml:"NODE"`
				LASTCHANGEDDATETIME string `xml:"LAST_CHANGED_DATE_TIME"`
				OBJECTCATEGORY      string `xml:"OBJECT_CATEGORY"`
				TECHCATEGORY        string `xml:"TECH_CATEGORY"`
				LIFECYCLESTAT       string `xml:"LIFE_CYCLE_STAT"`
				DUNAME              string `xml:"DU_NAME"`
				WRITEACCESS         string `xml:"WRITE_ACCESS"`
				DEPRECATED          string `xml:"DEPRECATED"`
				TRANSITIVEHASH      string `xml:"TRANSITIVE_HASH"`
				OFFLINEENABLED      string `xml:"OFFLINE_ENABLED"`
			} `xml:"item"`
		} `xml:"RESULT_SET"`
	} `xml:"values"`
}

// BCInformation struct
type BCInformation struct {
	XMLName xml.Name `xml:"abap"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Asx     string   `xml:"asx,attr"`
	Values  struct {
		Text      string `xml:",chardata"`
		RESULTSET struct {
			Text          string `xml:",chardata"`
			RELEASE       string `xml:"RELEASE"`
			DATATYPESTYPE struct {
				Text string `xml:",chardata"`
				Item struct {
					Text string `xml:",chardata"`
					// NAMESPACE
					// <NAMESPACE>http://0007042773-one-off.sap.com/Y7NLADCSY_</NAMESPACE>
					NAMESPACE string `xml:"NAMESPACE"`
					DATATYPES struct {
						Text string `xml:",chardata"`
						// <item>
						// 	<NODE_ID>00163E703A9E1EE997B7D4A29DBA98DC</NODE_ID>
						// 	<NAME>BCO_1ElementsQueryElements</NAME>
						// 	<PROXYNAME>Y7NLADCSY_BCTS720E0DE0CA1EB53D</PROXYNAME>
						// 	<LAST_CHANGED_DATE_TIME>2019-04-13T07:35:02Z</LAST_CHANGED_DATE_TIME>
						// 	<BASE_DT_KEY_NAME/>
						// 	<USAGE_CATEGORY>6</USAGE_CATEGORY>
						// 	<EXTENSIBLE/>
						// 	<REP_TERM>26</REP_TERM>
						// 	<TRANSITIVE_HASH>KSDjTUkBnkWvoVJhggkzgg==</TRANSITIVE_HASH>
						// </item>
						Item []struct {
							Text   string `xml:",chardata"`
							NODEID string `xml:"NODE_ID"`
							NAME   string `xml:"NAME"`
							// PROXYNAME, like
							PROXYNAME           string `xml:"PROXYNAME"`
							LASTCHANGEDDATETIME string `xml:"LAST_CHANGED_DATE_TIME"`
							BASEDTKEYNAME       string `xml:"BASE_DT_KEY_NAME"`
							USAGECATEGORY       string `xml:"USAGE_CATEGORY"`
							EXTENSIBLE          string `xml:"EXTENSIBLE"`
							REPTERM             string `xml:"REP_TERM"`
							TRANSITIVEHASH      string `xml:"TRANSITIVE_HASH"`
						} `xml:"item"`
					} `xml:"DATATYPES"`
				} `xml:"item"`
			} `xml:"DATATYPESTYPE"`
		} `xml:"RESULT_SET"`
	} `xml:"values"`
}

// BCPartnerSolution struct
type BCPartnerSolution struct {
	XMLName xml.Name `xml:"BCPartnerSolution"`
	Text    string   `xml:",chardata"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Head    struct {
		Text        string `xml:",chardata"`
		Xmlns       string `xml:"xmlns,attr"`
		ElementID   string `xml:"ElementID"`
		ElementType string `xml:"ElementType"`
	} `xml:"Head"`
	Bac struct {
		Text                      string `xml:",chardata"`
		Xmlns                     string `xml:"xmlns,attr"`
		ElementID                 string `xml:"ElementID"`
		ParentID                  string `xml:"ParentID"`
		VisibleFineTuning         string `xml:"VisibleFineTuning"`
		Description               string `xml:"Description"`
		ScopingQuestion           string `xml:"ScopingQuestion"`
		Type                      string `xml:"Type"`
		GoLiveActivityDescription string `xml:"GoLiveActivityDescription"`
		KTOverview                string `xml:"KTOverview"`
		KTRelevance               string `xml:"KTRelevance"`
	} `xml:"Bac"`
	Content []struct {
		Text        string `xml:",chardata"`
		Xmlns       string `xml:"xmlns,attr"`
		ID          string `xml:"ID"`
		Type        string `xml:"Type"`
		Description string `xml:"Description"`
		ElementID   string `xml:"ElementID"`
		Internal    string `xml:"Internal"`
	} `xml:"Content"`
}

// OberonApplication config
type OberonApplication struct {
	XMLName                         xml.Name `xml:"OberonApplication"`
	Text                            string   `xml:",chardata"`
	Xmlns                           string   `xml:"xmlns,attr"`
	StartComponent                  string   `xml:"startComponent,attr"`
	AllowTabbedBrowsing             string   `xml:"allowTabbedBrowsing,attr"`
	SuppressTabbedBrowsingControl   string   `xml:"suppressTabbedBrowsingControl,attr"`
	SuppressControlCenterHome       string   `xml:"suppressControlCenterHome,attr"`
	TextBundleResourceUriPrefix     string   `xml:"textBundleResourceUriPrefix,attr"`
	SuppressSystemMessages          string   `xml:"suppressSystemMessages,attr"`
	CapitalizeHeaderTexts           string   `xml:"capitalizeHeaderTexts,attr"`
	DefaultFeedsSearchCategory      string   `xml:"defaultFeedsSearchCategory,attr"`
	DefaultAddresseesSearchCategory string   `xml:"defaultAddresseesSearchCategory,attr"`
	RightAlignedLabels              string   `xml:"rightAlignedLabels,attr"`
	BasedOnBuiltInTheme             string   `xml:"basedOnBuiltInTheme,attr"`
	HasListSelectionColumn          string   `xml:"hasListSelectionColumn,attr"`
	UseBreadcrumbNavigationTitles   string   `xml:"useBreadcrumbNavigationTitles,attr"`
	ThemeUri                        string   `xml:"themeUri,attr"`
	ThemeHighContrastBlackUri       string   `xml:"themeHighContrastBlackUri,attr"`
	IsAutomationActive              string   `xml:"isAutomationActive,attr"`
	SolutionInfo                    struct {
		Text            string `xml:",chardata"`
		AssemblyVersion string `xml:"assemblyVersion,attr"`
		Changelist      string `xml:"changelist,attr"`
		Timestamp       string `xml:"timestamp,attr"`
		Codeline        string `xml:"codeline,attr"`
	} `xml:"SolutionInfo"`
	ProductSettings struct {
		Text                        string `xml:",chardata"`
		TitlePropertyName           string `xml:"titlePropertyName,attr"`
		AboutDialogImageResourceUri string `xml:"aboutDialogImageResourceUri,attr"`
		LogoResourceUri             string `xml:"logoResourceUri,attr"`
		LogoWhiteResourceUri        string `xml:"logoWhiteResourceUri,attr"`
		LogoutImageResourceUri      string `xml:"logoutImageResourceUri,attr"`
		SmartPhoneLoginLogo         string `xml:"smartPhoneLoginLogo,attr"`
		SmartPhoneLoginImage        string `xml:"smartPhoneLoginImage,attr"`
		SAPStandardText             string `xml:"SAPStandardText,attr"`
	} `xml:"ProductSettings"`
	LoginDialog struct {
		Text      string `xml:",chardata"`
		TitleText string `xml:"titleText,attr"`
		ImageUrl  string `xml:"imageUrl,attr"`
	} `xml:"LoginDialog"`
	SideCarPanels struct {
		Text                   string `xml:",chardata"`
		InitialWidth           string `xml:"initialWidth,attr"`
		ControllerSideCarPanel []struct {
			Text                      string `xml:",chardata"`
			ID                        string `xml:"id,attr"`
			ComponentName             string `xml:"componentName,attr"`
			SortIndex                 string `xml:"sortIndex,attr"`
			StartupTitleTextBundleKey string `xml:"startupTitleTextBundleKey,attr"`
			SupportOnly               string `xml:"supportOnly,attr"`
		} `xml:"ControllerSideCarPanel"`
		BuiltinSideCarPanel []struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id,attr"`
			Type      string `xml:"type,attr"`
			SortIndex string `xml:"sortIndex,attr"`
			IsModal   string `xml:"isModal,attr"`
		} `xml:"BuiltinSideCarPanel"`
	} `xml:"SideCarPanels"`
	ToolPaletteItems struct {
		Text                   string `xml:",chardata"`
		BuiltinToolPaletteItem []struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Type string `xml:"type,attr"`
		} `xml:"BuiltinToolPaletteItem"`
		NotificationCounterToolPaletteItem struct {
			Text       string `xml:",chardata"`
			ID         string `xml:"id,attr"`
			Icon       string `xml:"icon,attr"`
			TextPoolId string `xml:"textPoolId,attr"`
		} `xml:"NotificationCounterToolPaletteItem"`
		ToolPaletteSeparator struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"ToolPaletteSeparator"`
	} `xml:"ToolPaletteItems"`
	MenuItems struct {
		Text     string `xml:",chardata"`
		MenuItem []struct {
			Text               string `xml:",chardata"`
			ID                 string `xml:"id,attr"`
			TextKey            string `xml:"textKey,attr"`
			DefaultText        string `xml:"defaultText,attr"`
			ModeVisibility     string `xml:"modeVisibility,attr"`
			KeyUserOnly        string `xml:"keyUserOnly,attr"`
			NavigationMenuItem struct {
				Text            string `xml:",chardata"`
				ID              string `xml:"id,attr"`
				ComponentName   string `xml:"componentName,attr"`
				StartupPortName string `xml:"startupPortName,attr"`
				DefaultText     string `xml:"defaultText,attr"`
				TextKey         string `xml:"textKey,attr"`
			} `xml:"NavigationMenuItem"`
			BuiltinMenuItem []struct {
				Text           string `xml:",chardata"`
				ID             string `xml:"id,attr"`
				Type           string `xml:"type,attr"`
				DefaultText    string `xml:"defaultText,attr"`
				ModeVisibility string `xml:"modeVisibility,attr"`
				SupportOnly    string `xml:"supportOnly,attr"`
			} `xml:"BuiltinMenuItem"`
			BackgroundImageMenuItem struct {
				Text        string `xml:",chardata"`
				ID          string `xml:"id,attr"`
				DefaultText string `xml:"defaultText,attr"`
				KeyUser     string `xml:"keyUser,attr"`
			} `xml:"BackgroundImageMenuItem"`
			Separator []struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"Separator"`
			ExternalLinkMenuItem []struct {
				Text                   string `xml:",chardata"`
				ID                     string `xml:"id,attr"`
				Type                   string `xml:"type,attr"`
				CalculatedResourceName string `xml:"calculatedResourceName,attr"`
				Target                 string `xml:"target,attr"`
				DefaultText            string `xml:"defaultText,attr"`
				TextKey                string `xml:"textKey,attr"`
			} `xml:"ExternalLinkMenuItem"`
		} `xml:"MenuItem"`
		Separator []struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"Separator"`
		BuiltinMenuItem struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"id,attr"`
			Type        string `xml:"type,attr"`
			DefaultText string `xml:"defaultText,attr"`
		} `xml:"BuiltinMenuItem"`
		NavigationMenuItem struct {
			Text          string `xml:",chardata"`
			ID            string `xml:"id,attr"`
			ComponentName string `xml:"componentName,attr"`
			DefaultText   string `xml:"defaultText,attr"`
			TextKey       string `xml:"textKey,attr"`
		} `xml:"NavigationMenuItem"`
	} `xml:"MenuItems"`
	NotificationSettings struct {
		Text                       string `xml:",chardata"`
		NotificationCounterSetting struct {
			Text                           string `xml:",chardata"`
			ID                             string `xml:"id,attr"`
			NotificationCounterSettingItem struct {
				Text                   string `xml:",chardata"`
				ID                     string `xml:"id,attr"`
				TargetComponent        string `xml:"targetComponent,attr"`
				BTMNotificationCounter []struct {
					Text          string `xml:",chardata"`
					BtmTaskTypeID string `xml:"btmTaskTypeID,attr"`
				} `xml:"BTMNotificationCounter"`
			} `xml:"NotificationCounterSettingItem"`
		} `xml:"NotificationCounterSetting"`
	} `xml:"NotificationSettings"`
	ShellSettings struct {
		Text              string `xml:",chardata"`
		HideStartMenu     string `xml:"hideStartMenu,attr"`
		HideCommonTasks   string `xml:"hideCommonTasks,attr"`
		HideSingleHomeTab string `xml:"hideSingleHomeTab,attr"`
	} `xml:"ShellSettings"`
	Resources struct {
		Text               string `xml:",chardata"`
		CalculatedResource struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"CalculatedResource"`
	} `xml:"Resources"`
	ApplicationIcons struct {
		Text string `xml:",chardata"`
		Icon []struct {
			Text              string `xml:",chardata"`
			SymbolicName      string `xml:"symbolicName,attr"`
			ResourceUri       string `xml:"resourceUri,attr"`
			Width             string `xml:"width,attr"`
			Height            string `xml:"height,attr"`
			DedicatedForTheme string `xml:"dedicatedForTheme,attr"`
			HoverResourceUri  string `xml:"hoverResourceUri,attr"`
		} `xml:"Icon"`
	} `xml:"ApplicationIcons"`
}
