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
	listWidth  = 20
	listHeight = 14
	/*
	 * Fixed sizes for the list
	 * When non zero:
	 *   Width and or height are respecitvly set to listFixedWidth/listFixedHeight
	 * When 0:
	 *   Uses the screen width and/or height as listWidth/Height respectivly
	 */
	listFixedWidth  uint = 80
	listFixedHeight uint = 14
)

// styles
var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	checkStyle        = lipgloss.NewStyle().Background(lipgloss.ANSIColor(12))
)

const newTaskText string = "New Task"

// Repersents the state in a list
type ChecklistState int

const (
	// State for viewing the list, ie the default state
	StateList ChecklistState = iota
	// State for editing a task in the checklist
	StateEdit
	// State for adding a task into the checklist
	StateAdd
)

// Model that repersents the main state
type ChecklistModel struct {
	list  *list.Model   // Choosable items
	edit  EditTaskModel // Model to change
	state ChecklistState
}

// Default
func InitialModel(items []list.Item) ChecklistModel {
	list := list.New(items, itemDelegate{}, listWidth, listHeight)
	return ChecklistModel{
		list:  &list,
		edit:  EditTaskModel{}.New(),
		state: StateList,
	}
}

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

// Update a sub model of checklist model
func (m *ChecklistModel) UpdateSubModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var subModel tea.Model
	var cmd tea.Cmd = nil

	switch m.state {
	case StateList:
		var list list.Model
		list, cmd = m.list.Update(msg)
		if cmd != nil {
			return m, cmd
		}
		*(m.list) = list
		break
	case StateEdit:
		subModel, cmd = m.edit.Update(msg)
		m.edit = subModel.(EditTaskModel)
		break
	}

	return m, cmd
}

// Perform any closing actions needed for the current state
func (m ChecklistModel) ExitState() {
	switch m.state {
	case StateEdit:
		// Set the original item to the edited item
		m.list.Items()[m.edit.EditingItemIndex] = *m.edit.Item
		break
	}
}

func (m ChecklistModel) UpdateStateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			if m.state == StateList {
				selectedItem := m.list.SelectedItem()
				selectedCheckbox := selectedItem.(ChecklistItem)
				selectedCheckbox.Toggle()

				m.list.Items()[m.list.Index()] = selectedCheckbox
			}
			break
		case "e", "ctrl+e":
			// Exit current state
			m.ExitState()
			// Enter edit state
			// Set the pointer to the new edit one, store the index of the item
			item := m.list.Items()[m.list.Index()].(ChecklistItem)
			m.edit.SetItem(&item, m.list.Index())
			m.state = StateEdit
			return m, nil
		case "a", "ctrl+a":
			// item :=
			// var count = len(m.list.Items())
			cmd := m.list.InsertItem(m.list.Index(), ChecklistItem{newTaskText, false})
			// fmt.Printf("Size is %d, should be %d", len(m.list.Items()), count+1)
			return m, cmd
		case "d", "ctrl+d":
			m.list.RemoveItem(m.list.Index())
			return m, nil
		default:
			_, cmd := m.UpdateSubModel(msg)
			// m = mod.(ChecklistModel)
			return m, cmd
		}
	default:
		_, cmd := m.UpdateSubModel(msg)
		// m = mod.(ChecklistModel)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil

}

// Update the main checklist model
func (m ChecklistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check if we should quit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// TODO: Find way to update the list properly
			// WriteChecklist(OutPath, m)
			fmt.Printf("Exit Size is %d\n", len(m.list.Items()))
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Only do this in the list state to avoid race conditions
		if m.state == StateList {
			// Update width and height. See comments at listFixedWith for detail
			if listFixedWidth == 0 {
				listWidth = msg.Width
			} else {
				listWidth = int(listFixedWidth)
			}
			if listFixedHeight == 0 {
				listHeight = msg.Height
			} else {
				listHeight = int(listFixedHeight)
			}
			m.list.SetSize(listWidth, listHeight)
		}
	}

	// Pass update to proper update function
	switch m.state {
	case StateList:
		return m.UpdateStateList(msg)
	case StateEdit:
		return m.UpdateStateEdit(msg)
	}

	return m, nil
}

func (m ChecklistModel) View() string {
	// Send the UI for rendering
	switch m.state {
	case StateList:
		return m.list.View()
	case StateEdit:
		return m.edit.View()
	}

	return ""
}
