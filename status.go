package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusError
)

type Status struct {
	message    string
	statusType StatusType
	width      int
}

func NewStatus() Status {
	return Status{
		message:    "Ready. Type a command to execute.",
		statusType: StatusInfo,
	}
}

func (s *Status) SetWidth(width int) {
	s.width = width
}

func (s *Status) SetMessage(msg string, st StatusType) {
	s.message = msg
	s.statusType = st
}

func (s *Status) SetSuccess(msg string) {
	s.SetMessage(msg, StatusSuccess)
}

func (s *Status) SetError(msg string) {
	s.SetMessage(msg, StatusError)
}

func (s *Status) SetInfo(msg string) {
	s.SetMessage(msg, StatusInfo)
}

func (s Status) View() string {
	var style lipgloss.Style
	switch s.statusType {
	case StatusSuccess:
		style = StatusSuccessStyle
	case StatusError:
		style = StatusErrorStyle
	default:
		style = StatusInfoStyle
	}

	// Flatten to single line and truncate if too long
	msg := strings.ReplaceAll(s.message, "\n", " ")
	msg = strings.ReplaceAll(msg, "\r", "")
	msg = strings.TrimSpace(msg)
	maxLen := s.width - 4
	if maxLen > 0 && len(msg) > maxLen {
		msg = msg[:maxLen-3] + "..."
	}

	return style.Width(s.width).Render(msg)
}
