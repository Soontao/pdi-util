package pdiutil

import (
	"github.com/Soontao/pdi-util/ast"
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
	// not used now
	BOFieldsCount              int64
	CommunicationScenarioCount int64
	CodeListCount              int64
	BCOCount                   int64
	InternalCommunicationCount int64
}

// CountElementForBODL type
// if parse failed, return zero
func CountElementForBODL(source []byte) int {
	rt := 0

	if n, err := ast.ParseAST(source); err == nil && n != nil {
		if bo := n.GetNode("BODefinition"); bo != nil {
			if elements := bo.GetNodeList("Elements"); elements != nil {
				for _, e := range elements {
					switch e.GetType() {
					case "ElementItem":
						rt++
					}
				}
			}
		}
	}

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
		case "pid":
			rt.InternalCommunicationCount++
		}

	}

	for _, abslCode := range c.fetchSources(abslList, concurrent) {
		rt.ABSLCodeLines += len(strings.Split(abslCode.String(), "\n"))
	}

	return rt
}
