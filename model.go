package main

import (
	"fmt"
	"os"
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
	stdinData            []byte // piped input to feed to commands
	outputFile           string // file to write output to
	inputMode            string // "", "output_file"
	inputBuffer          string // buffer for input mode
}

func New(stdinData []byte, outputFile string) Model {
	output := NewOutput()
	status := NewStatus()
	// Show stdin data initially if present
	if len(stdinData) > 0 {
		output.SetContent(string(stdinData))
		// Write to output file if configured
		if outputFile != "" {
			if err := os.WriteFile(outputFile, stdinData, 0644); err != nil {
				status.SetError(fmt.Sprintf("Failed to write: %s", err))
			} else {
				status.SetSuccess(fmt.Sprintf("Wrote to %s", outputFile))
			}
		}
	}
	return Model{
		editor:     NewEditor(),
		output:     output,
		status:     status,
		stdinData:  stdinData,
		outputFile: outputFile,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle input mode (for setting output file)
	if m.inputMode != "" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.inputMode == "output_file" {
					m.outputFile = m.inputBuffer
					if m.outputFile != "" {
						// If editor empty and we have stdin, write it now
						if m.editor.Value() == "" && len(m.stdinData) > 0 {
							if err := os.WriteFile(m.outputFile, m.stdinData, 0644); err != nil {
								m.status.SetError(fmt.Sprintf("Failed to write: %s", err))
							} else {
								m.status.SetSuccess(fmt.Sprintf("Wrote to %s", m.outputFile))
							}
						} else {
							m.status.SetInfo(fmt.Sprintf("Output file set: %s", m.outputFile))
						}
					} else {
						m.status.SetInfo("Output file cleared")
					}
				}
				m.inputMode = ""
				m.inputBuffer = ""
				return m, nil
			case "esc":
				m.inputMode = ""
				m.inputBuffer = ""
				m.status.SetInfo("Cancelled")
				return m, nil
			case "backspace":
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.inputBuffer += msg.String()
				}
				return m, nil
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+d":
			m.editor.Reset()
			m.pendingCommand = ""
			// Show stdin when editor is cleared and write to output
			if len(m.stdinData) > 0 {
				m.output.SetContent(string(m.stdinData))
				m.output.GotoTop()
				if m.outputFile != "" {
					if err := os.WriteFile(m.outputFile, m.stdinData, 0644); err != nil {
						m.status.SetError(fmt.Sprintf("Cleared, but failed to write: %s", err))
						return m, nil
					}
					m.status.SetInfo(fmt.Sprintf("Editor cleared, wrote to %s", m.outputFile))
					return m, nil
				}
			}
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
		case "ctrl+o":
			m.inputMode = "output_file"
			m.inputBuffer = m.outputFile
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
			return m, Execute(currentCommand, m.stdinData)
		}
		return m, nil

	case CommandResultMsg:
		m.executing = false
		if msg.ExitCode == 0 {
			m.lastSuccessfulOutput = msg.Output
			m.lastCommand = m.editor.Value()
			m.output.SetContent(msg.Output)
			m.output.GotoTop()
			// Write to output file if configured
			if m.outputFile != "" {
				if err := os.WriteFile(m.outputFile, []byte(msg.Output), 0644); err != nil {
					m.status.SetError(fmt.Sprintf("✓ Command succeeded but failed to write: %s", err))
					return m, nil
				}
				m.status.SetSuccess(fmt.Sprintf("✓ Wrote to %s", m.outputFile))
			} else {
				m.status.SetSuccess("✓ Command succeeded")
			}
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

	// Show input prompt if in input mode
	if m.inputMode == "output_file" {
		prompt := fmt.Sprintf("Output file: %s█", m.inputBuffer)
		help := HelpStyle.Render("enter: confirm • esc: cancel")
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.editor.View(),
			m.output.View(),
			prompt,
			help,
		)
	}

	freezeStatus := ""
	if m.frozen {
		freezeStatus = " [❄ FROZEN]"
	}
	outputStatus := ""
	if m.outputFile != "" {
		outputStatus = fmt.Sprintf(" [-> %s]", m.outputFile)
	}
	help := HelpStyle.Render("^C: quit • ^D: clear • ^L: clear out • ^F: freeze • ^O: output file • pgup/dn: scroll" + freezeStatus + outputStatus)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.editor.View(),
		m.output.View(),
		m.status.View(),
		help,
	)
}
