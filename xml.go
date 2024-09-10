package main

import (
	"encoding/xml"
)

type XMLItemModel struct {
	XMLName xml.Name `xml:"task"`
	Title   string   `xml:"title,attr"`
	Checked bool     `xml:"checked,attr"`
}

func (m ChecklistItem) EncodeChecklistItem() XMLItemModel {
	return XMLItemModel{
		Title:   m.title,
		Checked: m.checked,
	}
}

func (m ChecklistModel) EncodeChecklist() (buf []byte, err error) {
	buf = make([]byte, 0)
	err = nil

	// Encode each checklist item
	type XMLChecklistModel struct {
		XMLName xml.Name        `xml:"checklist"`
		Items   []*XMLItemModel `xml:">task"`
	}

	model := XMLChecklistModel{
		Items: make([]*XMLItemModel, 0),
	}

	// Add items to xml model
	for _, v := range m.list.Items() {
		check := v.(ChecklistItem)
		item := check.EncodeChecklistItem()
		model.Items = append(model.Items, &item)
	}

	return xml.MarshalIndent(model, "", "\t")
}
