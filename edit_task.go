package main

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
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		}
	}
	return m, nil
}

func (m EditTaskModel) View() string {
	str := fmt.Sprintln(m.Item.Title)
	return str
}
