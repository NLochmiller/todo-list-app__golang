package main

import (
	"encoding/xml"

	"github.com/charmbracelet/bubbles/list"
)

type XMLItemModel struct {
	XMLName xml.Name `xml:"task"`
	Title   string   `xml:"title,attr"`
	Checked bool     `xml:"checked,attr"`
}

type XMLChecklistModel struct {
	XMLName xml.Name        `xml:"checklist"`
	Items   []*XMLItemModel `xml:">task"`
}

func (m ChecklistItem) EncodeChecklistItem() XMLItemModel {
	return XMLItemModel{
		Title:   m.Title,
		Checked: m.checked,
	}
}

func (m ChecklistModel) EncodeChecklist() (buf []byte, err error) {
	buf = make([]byte, 0)
	err = nil

	// Encode each checklist item

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

// Decode a single check item
func (m XMLItemModel) Decode() ChecklistItem {
	return ChecklistItem{
		Title:   m.Title,
		checked: m.Checked,
	}
}

// Decode an entire checklist from xml
func DecodeChecklist(buf []byte) (ChecklistModel, error) {
	// Parse xml into model
	var model XMLChecklistModel = XMLChecklistModel{}
	var err error = xml.Unmarshal(buf, &model)
	if err != nil {
		return ChecklistModel{}, err
	}

	// Parse xml model into checklist model
	var outputItems []list.Item = make([]list.Item, 0)
	for _, v := range model.Items {
		outputItems = append(outputItems, v.Decode())
	}
	return InitialModel(outputItems), err
}
