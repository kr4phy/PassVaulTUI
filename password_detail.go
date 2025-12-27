package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func UpdatePasswordDetail(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "enter":
			m.currentState = statePasswordsList
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func PasswordDetailView(m model) string {
	var details string
	if cred, ok := m.chosenCredential.(item); ok {
		details = fmt.Sprintf("Title: %s\nID: %s\nPassword: %s", keywordStyle.Render(cred.title[7:]), keywordStyle.Render(cred.id[4:]), keywordStyle.Render(cred.password))
	} else {
		details = "No credential selected."
	}

	content := fmt.Sprintf(
		"Credential Details\n\n%s\n\n(esc to go back)",
		details,
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
