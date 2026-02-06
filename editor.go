package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Editor struct {
	textarea textarea.Model
	width    int
	height   int
	focused  bool
}

func NewEditor() Editor {
	ta := textarea.New()
	ta.Placeholder = "Enter shell command..."
	ta.ShowLineNumbers = false
	ta.Prompt = "$ "
	ta.Focus()
	ta.CharLimit = 0 // No limit

	return Editor{
		textarea: ta,
		focused:  true,
	}
}

func (e *Editor) SetSize(width, height int) {
	e.width = width
	e.height = height
	// Account for border (2) and padding (2)
	e.textarea.SetWidth(width - 4)
	e.textarea.SetHeight(height - 2)
}

func (e *Editor) Focus() {
	e.focused = true
	e.textarea.Focus()
}

func (e *Editor) Blur() {
	e.focused = false
	e.textarea.Blur()
}

func (e *Editor) Value() string {
	return e.textarea.Value()
}

func (e *Editor) SetValue(s string) {
	e.textarea.SetValue(s)
}

func (e *Editor) Reset() {
	e.textarea.Reset()
}

func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	var cmd tea.Cmd
	e.textarea, cmd = e.textarea.Update(msg)
	return e, cmd
}

func (e Editor) View() string {
	style := EditorStyle
	if e.focused {
		style = EditorFocusedStyle
	}
	return style.Width(e.width - 2).Height(e.height - 2).Render(e.textarea.View())
}
