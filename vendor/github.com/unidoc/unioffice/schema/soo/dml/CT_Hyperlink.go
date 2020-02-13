// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package dml

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/unidoc/unioffice"
)

type CT_Hyperlink struct {
	IdAttr             *string
	InvalidUrlAttr     *string
	ActionAttr         *string
	TgtFrameAttr       *string
	TooltipAttr        *string
	HistoryAttr        *bool
	HighlightClickAttr *bool
	EndSndAttr         *bool
	Snd                *CT_EmbeddedWAVAudioFile
	ExtLst             *CT_OfficeArtExtensionList
}

func NewCT_Hyperlink() *CT_Hyperlink {
	ret := &CT_Hyperlink{}
	return ret
}

func (m *CT_Hyperlink) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.IdAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:id"},
			Value: fmt.Sprintf("%v", *m.IdAttr)})
	}
	if m.InvalidUrlAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "invalidUrl"},
			Value: fmt.Sprintf("%v", *m.InvalidUrlAttr)})
	}
	if m.ActionAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "action"},
			Value: fmt.Sprintf("%v", *m.ActionAttr)})
	}
	if m.TgtFrameAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "tgtFrame"},
			Value: fmt.Sprintf("%v", *m.TgtFrameAttr)})
	}
	if m.TooltipAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "tooltip"},
			Value: fmt.Sprintf("%v", *m.TooltipAttr)})
	}
	if m.HistoryAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "history"},
			Value: fmt.Sprintf("%d", b2i(*m.HistoryAttr))})
	}
	if m.HighlightClickAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "highlightClick"},
			Value: fmt.Sprintf("%d", b2i(*m.HighlightClickAttr))})
	}
	if m.EndSndAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "endSnd"},
			Value: fmt.Sprintf("%d", b2i(*m.EndSndAttr))})
	}
	e.EncodeToken(start)
	if m.Snd != nil {
		sesnd := xml.StartElement{Name: xml.Name{Local: "a:snd"}}
		e.EncodeElement(m.Snd, sesnd)
	}
	if m.ExtLst != nil {
		seextLst := xml.StartElement{Name: xml.Name{Local: "a:extLst"}}
		e.EncodeElement(m.ExtLst, seextLst)
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_Hyperlink) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "id" ||
			attr.Name.Space == "http://purl.oclc.org/ooxml/officeDocument/relationships" && attr.Name.Local == "id" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.IdAttr = &parsed
			continue
		}
		if attr.Name.Local == "invalidUrl" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.InvalidUrlAttr = &parsed
			continue
		}
		if attr.Name.Local == "action" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.ActionAttr = &parsed
			continue
		}
		if attr.Name.Local == "tgtFrame" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.TgtFrameAttr = &parsed
			continue
		}
		if attr.Name.Local == "tooltip" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.TooltipAttr = &parsed
			continue
		}
		if attr.Name.Local == "history" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.HistoryAttr = &parsed
			continue
		}
		if attr.Name.Local == "highlightClick" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.HighlightClickAttr = &parsed
			continue
		}
		if attr.Name.Local == "endSnd" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.EndSndAttr = &parsed
			continue
		}
	}
lCT_Hyperlink:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/main", Local: "snd"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/main", Local: "snd"}:
				m.Snd = NewCT_EmbeddedWAVAudioFile()
				if err := d.DecodeElement(m.Snd, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/main", Local: "extLst"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/main", Local: "extLst"}:
				m.ExtLst = NewCT_OfficeArtExtensionList()
				if err := d.DecodeElement(m.ExtLst, &el); err != nil {
					return err
				}
			default:
				unioffice.Log("skipping unsupported element on CT_Hyperlink %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lCT_Hyperlink
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the CT_Hyperlink and its children
func (m *CT_Hyperlink) Validate() error {
	return m.ValidateWithPath("CT_Hyperlink")
}

// ValidateWithPath validates the CT_Hyperlink and its children, prefixing error messages with path
func (m *CT_Hyperlink) ValidateWithPath(path string) error {
	if m.Snd != nil {
		if err := m.Snd.ValidateWithPath(path + "/Snd"); err != nil {
			return err
		}
	}
	if m.ExtLst != nil {
		if err := m.ExtLst.ValidateWithPath(path + "/ExtLst"); err != nil {
			return err
		}
	}
	return nil
}