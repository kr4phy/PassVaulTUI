package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var passListStyle = lipgloss.NewStyle().Padding(1, 2)

type item struct {
	title, id string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.id }
func (i item) FilterValue() string { return i.title }

func UpdatePasswordList(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := windowStyle.GetFrameSize()
		m.passList.SetSize(msg.Width-h, msg.Height-v-5)
	}

	var cmd tea.Cmd
	m.passList, cmd = m.passList.Update(msg)
	return m, cmd
}

func PasswordListView(m model) string {
	h, v := windowStyle.GetFrameSize()
	content := passListStyle.Render(m.passList.View())
	return lipgloss.Place(
		m.vpWidth-h,
		m.vpHeight-v,
		lipgloss.Center,
		lipgloss.Bottom,
		content,
	)
}
