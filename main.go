package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const DoDebug = true

func main() {
	choices := []list.Item{
		ChecklistItem{"XML encoding", true},
		ChecklistItem{"XML decoding", true},
		ChecklistItem{"Storage", false},
		ChecklistItem{"Better styling", false},
		ChecklistItem{"Editing", false},
	}

	var m ChecklistModel = InitialModel(choices)

	if DoDebug {
		var err error
		m, err = ReadChecklist("./input.xml")
		if err != nil {
			log.Fatal(err)
		}
	}

	m.list.SetShowTitle(false)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	/* Example
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
	*/

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

	WriteChecklist("./output.xml", m)
}
