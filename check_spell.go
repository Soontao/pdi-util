package pdiutil

import (
	"strings"

	"github.com/Jeffail/tunny"

	"github.com/beevik/etree"

	"github.com/tidwall/gjson"

	"github.com/imroc/req"
)

// DefaultRapidAPIToken test key, maybe deprecated
var DefaultRapidAPIToken = "92227c001emsh8847b5b6c9eced3p1163b0jsnc6a4f3a23aaa"

// UISuffixs list
var UISuffixs = []string{".uicomponent", ".uiwoc", ".uiwocview"}

// SpellErrorCheckResult type
type SpellErrorCheckResult struct {
	// File xrep file information
	File *XrepFile
	// These words maybe misspelling
	ErrorSpellingWords []string
}

// CheckSpellErrorAPI func
func (c *PDIClient) CheckSpellErrorAPI(solution, apiToken string, concurrent int) []*SpellErrorCheckResult {
	rt := []*SpellErrorCheckResult{}
	xrepPathes := c.GetSolutionXrepFileList(solution)
	uiXrepPathes := []string{}

	for _, p := range xrepPathes {
		for _, s := range UISuffixs {
			if strings.HasSuffix(p, s) {
				uiXrepPathes = append(uiXrepPathes, p)
				break
			}
		}
	}

	uiFiles := c.fetchSources(uiXrepPathes, DefaultDownloadConcurrnet)

	pool := tunny.NewFunc(concurrent, func(file interface{}) interface{} {
		// ignore error
		if words, _ := SpellCheckAPI(apiToken, FindTextPoolFromSource(file.(*XrepFile).Source)); len(words) > 0 {
			return &SpellErrorCheckResult{file.(*XrepFile), words}
		} else {
			return nil
		}
	})

	defer pool.Close()

	for _, file := range uiFiles {
		rt = append(rt, pool.Process(file).(*SpellErrorCheckResult))
	}

	return rt
}

// FindTextPoolFromSource func
//
func FindTextPoolFromSource(s []byte) []string {
	rt := []string{}
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(s)
	if err != nil {
		return rt
	}

	for _, e := range doc.SelectElement("UXComponent").SelectElement("TextPool").SelectElement("TextBlock").SelectElements("TextPoolEntry") {
		if t := e.SelectAttrValue("text", ""); t != "" {
			rt = append(rt, t)
		}
	}
	return rt
}

// SpellCheckAPI provided by Rapid API
//
// API Key required
func SpellCheckAPI(token string, words []string) ([]string, error) {
	rt := []string{}
	res, err := req.Get(
		"https://montanaflynn-spellcheck.p.rapidapi.com/check/",
		req.QueryParam{"text": strings.Join(words, " ")},
		req.Header{
			"X-RapidAPI-Host": "montanaflynn-spellcheck.p.rapidapi.com",
			"X-RapidAPI-Key":  token,
		},
	)

	if err != nil {
		return rt, err
	}

	gjson.Get(res.String(), "corrections").ForEach(func(key, value gjson.Result) bool {
		rt = append(rt, key.String())
		return true
	})

	return rt, nil
}
