package pdiutil

import (
	"strings"
	"time"
)

// ParseXrepDateString func
//
// input as '20181231092019.5268080' format
func ParseXrepDateString(input string) time.Time {
	sDateTime := strings.SplitN(input, ".", 2)[0]
	rt, _ := time.Parse("20060102030405", strings.TrimSpace(sDateTime))
	return rt
}
