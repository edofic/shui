package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
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
	// Handle init subcommand before flag parsing
	if len(os.Args) > 1 && os.Args[1] == "init" {
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
	}

	outputFile := flag.String("o", "", "output file to write results to")
	flag.Parse()

	// Check if stdin is a pipe (not a terminal)
	var stdinData []byte
	var opts []tea.ProgramOption

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		// Read piped input to feed to commands
		var err error
		stdinData, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}

		// Open /dev/tty for interactive input
		tty, err := os.Open("/dev/tty")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening /dev/tty: %v\n", err)
			os.Exit(1)
		}
		defer tty.Close()

		opts = append(opts, tea.WithInput(tty))
	}

	opts = append(opts, tea.WithAltScreen())

	p := tea.NewProgram(
		New(stdinData, *outputFile),
		opts...,
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
