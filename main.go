package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	choices := []list.Item{
		ChecklistItem{"Storage", false},
		ChecklistItem{"Better styling", false},
		ChecklistItem{"Editing", false},
	}

	var m ChecklistModel = InitialModel(choices)
	m.list.SetShowTitle(false)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	items := m.list.Items()
	// Check each item
	for _, v := range items {
		// check if selected
		checked := " "
		if v.(ChecklistItem).checked {
			checked = "x" // selected!
		}

		fmt.Printf("[%s] %q\n", checked, v)
	}

	// Distance
	fmt.Printf("\n\n\n")

	// Testing
	buf, _ := m.EncodeChecklist()
	fmt.Println(string(buf))

	fmt.Printf("\n\n\n")
	mod, err := DecodeChecklist(buf)

	if err != nil {
		fmt.Errorf("a", err)
	}

	for _, v := range mod.list.Items() {
		item := v.(ChecklistItem)
		fmt.Printf("%q %t\n", item.Title, item.Checked())
	}
}
