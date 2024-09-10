package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	title   string // The display name of this item
	checked bool   // Is this checked off?
}

func (i ChecklistItem) FilterValue() string { return i.title }

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
	// cfn := uncheckedboxStyle.Render
	if i.Checked() {
		checked = "x"
		// cfn = checkedboxStyle.Render
	}
	// checked = cfn("[" + checked + "]")
	checked = "[" + checked + "]"

	str := fmt.Sprintf("%d. %s", index+1, i.title)
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

		// // The "up" and "k" keys move the cursor up
		// case "up", "k":
		// 	if m.list.Cursor() > 0 {
		// 		m.list.CursorUp()
		// 	}

		// // The "down" and "j" keys move the cursor down
		// case "down", "j":
		// 	if m.list.Cursor() < len(m.list.Items())-1 {
		// 		m.list.CursorDown()
		// 	}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selectedItem := m.list.SelectedItem()
			selectedCheckbox := selectedItem.(ChecklistItem)
			selectedCheckbox.Toggle()
			m.list.Items()[m.list.Cursor()] = selectedCheckbox
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m ChecklistModel) View() string {
	// // The header
	// s := "What should we buy at the market?\n\n"

	// // Iterate over our choices
	// for i, choice := range m.list.Items() {

	// 	// Is the .list.Cursor() pointing at this choice?
	// 	.list.Cursor() := " " // no .list.Cursor()
	// 	if m..list.Cursor() == i {
	// 		.list.Cursor() = ">" // .list.Cursor()!
	// 	}

	// 	// Is this choice selected?
	// 	checked := " " // not selected
	// 	if _, ok := m.selected[i]; ok {
	// 		checked = "x" // selected!
	// 	}

	// 	// Render the row
	// 	s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	// }

	// // The footer
	// s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return "\n" + m.list.View()
}

type ChecklistModel struct {
	list     list.Model       // Choosable items
}

// Default
func InitialModel(items []list.Item) ChecklistModel {
	const defaultWidth = 20
	return ChecklistModel{
		list: list.New(items, itemDelegate{}, 20, 14),
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
	}
}
