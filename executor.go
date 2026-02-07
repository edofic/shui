package main

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultTimeout = 30 * time.Second

// Execute runs a shell command and returns a tea.Cmd that will produce a CommandResultMsg
func Execute(command string, stdinData []byte) tea.Cmd {
	return ExecuteWithTimeout(command, stdinData, defaultTimeout)
}

// ExecuteWithTimeout runs a shell command with a specified timeout
func ExecuteWithTimeout(command string, stdinData []byte, timeout time.Duration) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, "sh", "-c", command)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if len(stdinData) > 0 {
			cmd.Stdin = bytes.NewReader(stdinData)
		}

		err := cmd.Run()

		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				// Other error (e.g., command not found, timeout)
				return CommandResultMsg{
					Output:   "",
					Stderr:   err.Error(),
					ExitCode: -1,
					Err:      err,
				}
			}
		}

		return CommandResultMsg{
			Output:   stdout.String(),
			Stderr:   stderr.String(),
			ExitCode: exitCode,
			Err:      nil,
		}
	}
}
