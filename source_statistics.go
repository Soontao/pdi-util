package pdiutil

import (
	"github.com/Soontao/pdi-util/ast"
	"github.com/Soontao/pdi-util/ast/types"
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

// CountElementForBODL type
// if parse failed, return zero
func CountElementForBODL(source []byte) int {
	rt := 0

	if n, err := ast.ParseAST(source); err == nil && n != nil {
		if bo := n.GetNode("BODefinition"); bo != nil {
			rt = countNodeInnerElements(bo)
		}
	}

	return rt
}

func countNodeInnerElements(n *types.GrammerNode) int {
	rt := 0
	if elements := n.GetNodeList("Elements"); elements != nil {
		for _, e := range elements {
			switch e.GetType() {
			case "ElementItem":
				rt++
			case "BusinessObjectNode":
				rt += countNodeInnerElements(e)
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

	for _, boCode := range c.fetchSources(boList, concurrent) {
		rt.BOFieldsCount += CountElementForBODL(boCode.Source)
	}

	for _, uiCode := range c.fetchSources(uiList, concurrent) {
		rt.UIComplexity += CountXMLComplexity(uiCode.Source)
	}

	return rt
}
