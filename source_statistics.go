package pdiutil

import (
	"path/filepath"
	"strings"
)

// SolutionStatisticsResult type
type SolutionStatisticsResult struct {
	Solution                   *Solution
	ABSLCodeLines              int
	ABSLFileCount              int
	WebServicesCount           int
	UIComponentCount           int
	BOCount                    int
	BOFieldsCount              int
	CommunicationScenarioCount int
	CodeListCount              int
	BCOCount                   int
	InternalCommunicationCount int
	UIComplexity               int
}

// Statistics solution
func (c *PDIClient) Statistics(solution string, concurrent int) *SolutionStatisticsResult {
	s := c.GetSolutionByIDOrDescription(solution)
	rt := &SolutionStatisticsResult{Solution: &s}
	// retrieve files list
	files := c.GetSolutionXrepFileList(s.Name)

	abslList := []string{}

	boList := []string{}

	uiList := []string{}

	for _, f := range files {
		switch strings.TrimPrefix(filepath.Ext(f), ".") {
		case "uicomponent", "xuicomponent", "uiwoc", "uiwocview":
			rt.UIComponentCount++
			uiList = append(uiList, f)
		case "bo", "xbo":
			rt.BOCount++
			boList = append(boList, f)
		case "wsid":
			rt.WebServicesCount++
		case "absl":
			rt.ABSLFileCount++
			abslList = append(abslList, f)
		case "bco":
			rt.BCOCount++
		case "codelist":
			rt.CodeListCount++
		case "csd":
			rt.CommunicationScenarioCount++
		case "pid":
			rt.InternalCommunicationCount++
		}

	}

	for _, abslCode := range c.fetchSources(abslList, concurrent) {
		rt.ABSLCodeLines += len(strings.Split(abslCode.String(), "\n"))
	}

	return rt
}
