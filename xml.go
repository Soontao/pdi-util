package pdiutil

import (
	"github.com/clbanning/mxj"
)

// PrettifyXML func
//
// not work now
func PrettifyXML(xmls string) (rt string) {
	tmp, err := mxj.BeautifyXml([]byte(xmls), "", "  ")
	if err != nil {
		panic(err)
	} else {
		rt = string(tmp[:])
	}
	return rt
}
