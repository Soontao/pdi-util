package main

import (
	"github.com/clbanning/mxj"
)

func PrettifyXML(xmls string) (rt string) {
	tmp, err := mxj.BeautifyXml([]byte(xmls), "", "  ")
	if err != nil {
		panic(err)
	} else {
		rt = string(tmp[:])
	}
	return rt
}
