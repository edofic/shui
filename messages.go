package main

// DebounceTickMsg is sent after the debounce delay to trigger execution
type DebounceTickMsg struct {
	Command string // The command that was pending when the tick started
}

// CommandResultMsg contains the result of a shell command execution
type CommandResultMsg struct {
	Output   string
	Stderr   string
	ExitCode int
	Err      error
}
