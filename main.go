package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const zshInit = `shui() {
  local cmd
  cmd=$(command shui)
  if [[ -n "$cmd" ]]; then
    print -z "$cmd"
  fi
}
`

const bashInit = `shui() {
  local cmd
  cmd=$(command shui)
  if [[ -n "$cmd" ]]; then
    history -s "$cmd"
    bind '"\e[0n": "'"$cmd"'"'
    printf '\e[5n'
  fi
}
`

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: shui init <shell>")
				fmt.Fprintln(os.Stderr, "Supported shells: zsh, bash")
				os.Exit(1)
			}
			switch os.Args[2] {
			case "zsh":
				fmt.Print(zshInit)
			case "bash":
				fmt.Print(bashInit)
			default:
				fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", os.Args[2])
				fmt.Fprintln(os.Stderr, "Supported shells: zsh, bash")
				os.Exit(1)
			}
			return
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
			os.Exit(1)
		}
	}

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
