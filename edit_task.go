package main

// Do this
// https://github.com/charmbracelet/bubbletea/tree/master/examples/textinputs
import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

/* Implements the tea model interface */
type EditTaskModel struct {
	Item  *ChecklistItem
	Index int
}

func (m EditTaskModel) Init() tea.Cmd {
	return nil
}

func (m EditTaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// // Is it a key press?
	// case tea.KeyMsg:
	// 	switch msg.String() {
	// 	}
	// }
	return m, nil
}

func (m EditTaskModel) View() string {
	str := fmt.Sprintln(m.Item.Title)
	return str
}

// Handle update for when the main models state is StateEdit
func (m ChecklistModel) UpdateStateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "l":
			// Exit current state
			m.ExitState()
			// Enter list state
			m.state = StateList
			return m, nil
		default:
			return m.UpdateSubModel(msg)
		}
	default:
		return m.UpdateSubModel(msg)
	}
}
