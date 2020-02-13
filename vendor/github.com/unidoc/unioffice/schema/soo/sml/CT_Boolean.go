// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package sml

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/unidoc/unioffice"
)

type CT_Boolean struct {
	// Value
	VAttr bool
	// Unused Item
	UAttr *bool
	// Calculated Item
	FAttr *bool
	// Caption
	CAttr *string
	// Member Property Count
	CpAttr *uint32
	// Member Property Indexes
	X []*CT_X
}

func NewCT_Boolean() *CT_Boolean {
	ret := &CT_Boolean{}
	return ret
}

func (m *CT_Boolean) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "v"},
		Value: fmt.Sprintf("%d", b2i(m.VAttr))})
	if m.UAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "u"},
			Value: fmt.Sprintf("%d", b2i(*m.UAttr))})
	}
	if m.FAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "f"},
			Value: fmt.Sprintf("%d", b2i(*m.FAttr))})
	}
	if m.CAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "c"},
			Value: fmt.Sprintf("%v", *m.CAttr)})
	}
	if m.CpAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "cp"},
			Value: fmt.Sprintf("%v", *m.CpAttr)})
	}
	e.EncodeToken(start)
	if m.X != nil {
		sex := xml.StartElement{Name: xml.Name{Local: "ma:x"}}
		for _, c := range m.X {
			e.EncodeElement(c, sex)
		}
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_Boolean) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Local == "v" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.VAttr = parsed
			continue
		}
		if attr.Name.Local == "u" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.UAttr = &parsed
			continue
		}
		if attr.Name.Local == "f" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.FAttr = &parsed
			continue
		}
		if attr.Name.Local == "c" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.CAttr = &parsed
			continue
		}
		if attr.Name.Local == "cp" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.CpAttr = &pt
			continue
		}
	}
lCT_Boolean:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/spreadsheetml/2006/main", Local: "x"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/spreadsheetml/main", Local: "x"}:
				tmp := NewCT_X()
				if err := d.DecodeElement(tmp, &el); err != nil {
					return err
				}
				m.X = append(m.X, tmp)
			default:
				unioffice.Log("skipping unsupported element on CT_Boolean %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lCT_Boolean
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the CT_Boolean and its children
func (m *CT_Boolean) Validate() error {
	return m.ValidateWithPath("CT_Boolean")
}

// ValidateWithPath validates the CT_Boolean and its children, prefixing error messages with path
func (m *CT_Boolean) ValidateWithPath(path string) error {
	for i, v := range m.X {
		if err := v.ValidateWithPath(fmt.Sprintf("%s/X[%d]", path, i)); err != nil {
			return err
		}
	}
	return nil
}