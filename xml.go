package pdiutil

import (
	"bytes"
	"encoding/xml"
)

// XMLNode
type XMLNode struct {
	XMLName xml.Name
	Content []byte     `xml:",innerxml"`
	Nodes   []*XMLNode `xml:",any"`
}

// ParseXML to node
//
// Basic usage
func ParseXML(xmlBytes []byte) *XMLNode {
	rt := &XMLNode{}
	xml.NewDecoder(bytes.NewBuffer(xmlBytes)).Decode(rt)
	return rt
}

// WalkXMLNode func
//
// for callback function, if you dont want deep in specific node, just return false
func WalkXMLNode(nodes []*XMLNode, f func(*XMLNode) bool) {
	for _, n := range nodes {
		if f(n) {
			WalkXMLNode(n.Nodes, f)
		}
	}
}

// CountXMLComplexity number
func CountXMLComplexity(xmlBytes []byte) int {
	rt := 0
	WalkXMLNode([]*XMLNode{ParseXML(xmlBytes)}, func(n *XMLNode) bool {
		rt++
		return true
	})
	return rt
}
