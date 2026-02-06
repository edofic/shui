package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const debounceDelay = 300 * time.Millisecond

type Model struct {
	editor               Editor
	output               Output
	status               Status
	lastSuccessfulOutput string
	lastCommand          string
	executing            bool
	pendingCommand       string
	width, height        int
	ready                bool
	frozen               bool
}

func New() Model {
	return Model{
		editor: NewEditor(),
		output: NewOutput(),
		status: NewStatus(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+d":
			m.editor.Reset()
			m.pendingCommand = ""
			m.status.SetInfo("Editor cleared")
			return m, nil
		case "ctrl+l":
			m.output.SetContent("")
			m.lastSuccessfulOutput = ""
			m.status.SetInfo("Output cleared")
			return m, nil
		case "ctrl+f":
			m.frozen = !m.frozen
			if m.frozen {
				m.status.SetInfo("❄ Output frozen - editing won't trigger execution")
			} else {
				m.status.SetInfo("▶ Output unfrozen - resuming live updates")
				// Trigger execution of current command if there is one
				cmd := m.editor.Value()
				if cmd != "" && !m.executing {
					return m, m.startDebounce(cmd)
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateSizes()
		return m, nil

	case DebounceTickMsg:
		// Only execute if the command hasn't changed since debounce started
		// and we're not frozen
		currentCommand := m.editor.Value()
		if msg.Command == currentCommand && currentCommand != "" && !m.executing && !m.frozen {
			m.executing = true
			m.status.SetInfo("Executing...")
			return m, Execute(currentCommand)
		}
		return m, nil

	case CommandResultMsg:
		m.executing = false
		if msg.ExitCode == 0 {
			m.lastSuccessfulOutput = msg.Output
			m.lastCommand = m.editor.Value()
			m.output.SetContent(msg.Output)
			m.output.GotoTop()
			m.status.SetSuccess(fmt.Sprintf("✓ Command succeeded"))
		} else {
			// Keep previous successful output, show error in status
			errMsg := msg.Stderr
			if errMsg == "" && msg.Err != nil {
				errMsg = msg.Err.Error()
			}
			if errMsg == "" {
				errMsg = fmt.Sprintf("Exit code: %d", msg.ExitCode)
			}
			m.status.SetError(fmt.Sprintf("✗ %s", errMsg))
		}
		return m, nil
	}

	// Handle editor updates
	prevValue := m.editor.Value()
	var editorCmd tea.Cmd
	m.editor, editorCmd = m.editor.Update(msg)
	cmds = append(cmds, editorCmd)

	// Check if editor content changed - start debounce (only if not frozen)
	newValue := m.editor.Value()
	if newValue != prevValue {
		m.pendingCommand = newValue
		if newValue != "" && !m.frozen {
			cmds = append(cmds, m.startDebounce(newValue))
		}
	}

	// Handle output viewport updates (scrolling)
	var outputCmd tea.Cmd
	m.output, outputCmd = m.output.Update(msg)
	cmds = append(cmds, outputCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateSizes() {
	// Layout: editor on top (30% min 5 lines), output below, status bar at bottom
	editorHeight := m.height * 30 / 100
	if editorHeight < 5 {
		editorHeight = 5
	}
	if editorHeight > 10 {
		editorHeight = 10
	}

	statusHeight := 1
	outputHeight := m.height - editorHeight - statusHeight - 1 // -1 for help line

	m.editor.SetSize(m.width, editorHeight)
	m.output.SetSize(m.width, outputHeight)
	m.status.SetWidth(m.width)
}

func (m Model) startDebounce(command string) tea.Cmd {
	return tea.Tick(debounceDelay, func(_ time.Time) tea.Msg {
		return DebounceTickMsg{Command: command}
	})
}

func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	freezeStatus := ""
	if m.frozen {
		freezeStatus = " [❄ FROZEN]"
	}
	help := HelpStyle.Render("ctrl+c: quit • ctrl+d: clear editor • ctrl+l: clear output • ctrl+f: freeze/unfreeze • pgup/pgdn: scroll" + freezeStatus)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.editor.View(),
		m.output.View(),
		m.status.View(),
		help,
	)
}
