package main

import (
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

func (m EditTaskModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m EditTaskModel) View() string {
	return ""
}
