package pdiutil

import (
	"path/filepath"
	"strings"
)

// SolutionStatisticsResult type
type SolutionStatisticsResult struct {
	Solution         *Solution
	ABSLCodeLines    int
	ABSLFileCount    int64
	WebServicesCount int64
	UIComponentCount int64
	BOCount          int64
	// not use now
	BOFieldsCount              int64
	CommunicationScenarioCount int64
	CodeListCount              int64
	BCOCount                   int64
}

// CountElementForBODL type
func CountElementForBODL(bodl []byte) int {
	rt := 0

	return rt
}

// Statistics solution
func (c *PDIClient) Statistics(solution string, concurrent int) *SolutionStatisticsResult {
	s := c.GetSolutionByIDOrDescription(solution)
	rt := &SolutionStatisticsResult{Solution: &s}
	// retrieve files list
	files := c.GetSolutionXrepFileList(s.Name)

	abslList := []string{}

	boList := []string{}

	for _, f := range files {
		switch strings.TrimPrefix(filepath.Ext(f), ".") {
		case "uicomponent", "xuicomponent", "uiwoc", "uiwocview":
			rt.UIComponentCount++
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
		}
	}

	for _, abslCode := range c.fetchSources(abslList, concurrent) {
		rt.ABSLCodeLines += len(strings.Split(abslCode.String(), "\n"))
	}

	return rt
}
