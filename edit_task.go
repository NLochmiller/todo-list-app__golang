package main

// Do this
// https://github.com/charmbracelet/bubbletea/tree/master/examples/textinputs
import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style
var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

// How many input fields we want
const numInputs int = 1

/* Implements the tea model interface */
type EditTaskModel struct {
	// Variables related to the current model
	Item             *ChecklistItem
	EditingItemIndex int
	// Things needed for editing text
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func (m EditTaskModel) New() EditTaskModel {
	m = EditTaskModel{
		inputs: make([]textinput.Model, numInputs),
	}
	return m
}

// Set the item to given checklist item
func (m *EditTaskModel) SetItem(item *ChecklistItem, index int) {
	m.Item = item
	m.EditingItemIndex = index

	var t textinput.Model
	for i := range m.inputs {

		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 0 // Accept infinite characters

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.SetValue(m.Item.Title)
			// case 1:
			// 	t.Placeholder = "Email"
			// 	t.CharLimit = 64
			// case 2:
			// 	t.Placeholder = "Password"
			// 	t.EchoMode = textinput.EchoPassword
			// 	t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}
}

// Exits the edit state, does not save changes
func (m ChecklistModel) exitToList(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.ExitState()
	m.state = StateList
	return m, nil
}

/* Update handlers */
// Handle main update for when the main models state is StateEdit
func (m ChecklistModel) UpdateStateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s": // Save the items then return to the list state
			m.edit.saveInputsToItem()
			return m.exitToList(msg)
		case "ctrl+q": // Exit to the list state. Does not change current state
			return m.exitToList(msg)
		default:
			return m.UpdateSubModel(msg)
		}
	default:
		return m.UpdateSubModel(msg)
	}
}

// Update each text field in the editing model
func (m *EditTaskModel) updateInputFields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Get the data from the input into the item
func (m EditTaskModel) saveInputsToItem() {
	m.Item.Title = m.inputs[0].Value()
}

/* Implement tea.Model */
func (m EditTaskModel) Init() tea.Cmd {
	m = m.New()
	return textinput.Blink
}

func (m EditTaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				//TODO: May need to disable, need to find way to have
				// custom tea.Msg
				// m.saveInputsToItem()
				// return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputFields(msg)
	return m, cmd
}

func (m EditTaskModel) View() string {
	var b strings.Builder

	b.WriteString(helpStyle.Render("Editing task #"+strconv.Itoa(m.EditingItemIndex),
		string(m.Item.Title)))
	b.WriteRune('\n')

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
