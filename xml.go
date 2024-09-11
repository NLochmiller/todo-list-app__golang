/*
 * Convert checklists to xml, reading and writing that from the filesystem
 * Most users will want to use the ReadChecklist and WriteChecklist functions
 */
package main

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

// Read the xml file into a checklist, returns empty model if err != nil
func ReadChecklist(filepath string) (ChecklistModel, error) {
	buf, err := os.ReadFile(filepath)
	if err != nil && err != io.EOF {
		return ChecklistModel{}, err
	}
	if buf == nil {
		return ChecklistModel{}, os.ErrInvalid // TODO: Create custom error
	}
	return DecodeChecklist(buf)
}

// Wrapper for WriteChecklistPerms with hardcoded perms
func WriteChecklist(filepath string, model ChecklistModel) (err error) {
	return WriteChecklistPerms(filepath, model, 0644)
}

// Write the checklist converted to xml to the file at filepath creating it if
// the file does not exist
func WriteChecklistPerms(filepath string, model ChecklistModel, perm os.FileMode) (err error) {
	buf, err := model.EncodeChecklist()
	if err != nil {
		return
	}

	return os.WriteFile(filepath, buf, perm)
}

type XMLItemModel struct {
	XMLName xml.Name `xml:"task"`
	Title   string   `xml:"title,attr"`
	Checked bool     `xml:"checked,attr"`
}

type XMLChecklistModel struct {
	XMLName xml.Name        `xml:"checklist"`
	Items   []*XMLItemModel `xml:">task"`
}

func (m ChecklistItem) EncodeChecklistItem() *XMLItemModel {
	return &XMLItemModel{
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
		model.Items = append(model.Items, nil)
		model.Items[len(model.Items) - 1] =  check.EncodeChecklistItem()
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
