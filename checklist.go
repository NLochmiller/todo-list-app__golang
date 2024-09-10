package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listWidth  = 20
	listHeight = 14
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	checkStyle        = lipgloss.NewStyle().Background(lipgloss.ANSIColor(12))
)

/* Custom list item */
/* List ChecklistItem struct */
/* Tasks */
/* Implements list.Item */
type ChecklistItem struct {
	Title   string // The display name of this item
	checked bool   // Is this checked off?
}

func (i ChecklistItem) FilterValue() string { return i.Title }

func (i *ChecklistItem) SetChecked(b bool) {
	i.checked = b
}
func (i ChecklistItem) Checked() bool {
	return i.checked
}

// Simply invert checked
func (i *ChecklistItem) Toggle() {
	i.checked = !(i.checked)
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ChecklistItem)
	if !ok {
		return
	}

	// Create checkbox
	var checked string = " "
	if i.Checked() {
		checked = "x"
	}
	checked = "[" + checked + "]"

	str := fmt.Sprintf("%d. %s", index+1, i.Title)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(checked, str))
}

// Don't do any work for now
func (m ChecklistModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m ChecklistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	list, cmd := m.list.Update(msg)

	if cmd != nil {
		return m, cmd
	}
	m.list = list

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selectedItem := m.list.SelectedItem()
			selectedCheckbox := selectedItem.(ChecklistItem)
			selectedCheckbox.Toggle()
			// TODO: Needs to not use cursor but rather position in list
			m.list.Items()[m.list.Index()] = selectedCheckbox
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m ChecklistModel) View() string {
	// Send the UI for rendering
	return "\n" + m.list.View()
}

type ChecklistModel struct {
	list list.Model // Choosable items
}

// Default
func InitialModel(items []list.Item) ChecklistModel {
	return ChecklistModel{
		list: list.New(items, itemDelegate{}, listWidth, listHeight),
	}
}
