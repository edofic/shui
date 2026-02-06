package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
	)

	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	if model, ok := m.(Model); ok {
		if cmd := model.editor.Value(); cmd != "" {
			fmt.Println(cmd)
		}
	}
}
