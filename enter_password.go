package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func UpdateEnterPassword(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.masterPass = m.passwordInput.Value()
			m.currentState = statePasswordsList
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.passwordInput, cmd = m.passwordInput.Update(msg)
	return m, cmd
}

func EnterPasswordView(m model) string {
	content := fmt.Sprintf(
		"Please enter password.\n\n%s\n\n%s",
		m.passwordInput.View(),
		"(esc to quit)",
	) + "\n"
	h, v := windowStyle.GetFrameSize()
	return lipgloss.Place(
		m.vpWidth-h,
		m.vpHeight-v,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
