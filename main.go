package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}


	str := fmt.Sprintf("%d. %s", index+1, i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// Don't do any work for now
func (m checklistModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m checklistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.list.Cursor() > 0 {
				m.list.CursorUp()
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.list.Cursor() < len(m.list.Items())-1 {
				m.list.CursorDown()
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.list.Cursor()]
			if ok {
				delete(m.selected, m.list.Cursor())
			} else {
				m.selected[m.list.Cursor()] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m checklistModel) View() string {
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

type checklistModel struct {
	list     list.Model       // Choosable items
	selected map[int]struct{} // which to-do items are selected
}

// Default
func initialModel(items []list.Item) checklistModel {
	const defaultWidth = 20
	return checklistModel{
		list: list.New(items, itemDelegate{}, 20, 14),
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func main() {
	choices := []list.Item{
		item("Buy carrots"),
		item("Buy celery"),
		item("Buy kohlrabi")}

	var m checklistModel = initialModel(choices)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	items := m.list.Items()
	// Check each item
	for i, v := range items {
		// check if selected
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		fmt.Printf("[%s] %q\n", checked, v)
	}
}
