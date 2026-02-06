package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Output struct {
	viewport viewport.Model
	width    int
	height   int
	content  string
}

func NewOutput() Output {
	vp := viewport.New(0, 0)
	return Output{
		viewport: vp,
	}
}

func (o *Output) SetSize(width, height int) {
	o.width = width
	o.height = height
	// Account for border (2) and padding (2)
	o.viewport.Width = width - 4
	o.viewport.Height = height - 2
}

func (o *Output) SetContent(content string) {
	o.content = content
	o.viewport.SetContent(content)
}

func (o *Output) GotoTop() {
	o.viewport.GotoTop()
}

func (o *Output) GotoBottom() {
	o.viewport.GotoBottom()
}

func (o Output) Update(msg tea.Msg) (Output, tea.Cmd) {
	var cmd tea.Cmd
	o.viewport, cmd = o.viewport.Update(msg)
	return o, cmd
}

func (o Output) View() string {
	return OutputStyle.Width(o.width - 2).Height(o.height - 2).Render(o.viewport.View())
}

func (o Output) ScrollPercent() float64 {
	return o.viewport.ScrollPercent()
}
