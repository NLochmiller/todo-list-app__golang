package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const DoDebug = true

func GetExampleList() ChecklistModel {
	return InitialModel([]list.Item{
		ChecklistItem{"Fix multi page bug", true},
		ChecklistItem{"XML encoding", true},
		ChecklistItem{"XML decoding", true},
		ChecklistItem{"Storage", true},
		ChecklistItem{"User defined storage paths", false},
		ChecklistItem{"Editing", false},
		ChecklistItem{"Multistate boxes", false}, // [ ] [-] [!] [?] [x]
		ChecklistItem{"Better styling", false},
		ChecklistItem{"Placeholder 1", false},
		ChecklistItem{"Placeholder 2", false},
		ChecklistItem{"Placeholder 3", false},
		ChecklistItem{"Placeholder 4", false},
		ChecklistItem{"Placeholder 5", false},
		ChecklistItem{"Placeholder 6", false},
		ChecklistItem{"Placeholder 7", false},
		ChecklistItem{"Placeholder 8", false},
		ChecklistItem{"Placeholder 9", false},
		ChecklistItem{"Placeholder 10", false},
	})
}

var inPath, OutPath string = "../database.xml", "../database.xml"

func main() {
	var mi ChecklistModel = GetExampleList()
	var m *ChecklistModel = &mi

	// Load from input
	mod, err := ReadChecklist(inPath)
	// If there was an error that is not a child of ErrNotExist (ie no file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	} else if !errors.Is(err, os.ErrNotExist) {
		// Only override m if the file exists
		*m = mod
	}

	m.list.SetShowTitle(false)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	WriteChecklist(OutPath, *m)
}
