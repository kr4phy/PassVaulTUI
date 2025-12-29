package main

import (
	"errors"
	"fmt"
	"os"

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
			_, err := LoadEncryptedData(dataFilePath(), deriveKey(m.masterPass))

			if err != nil {
				switch {
				case errors.Is(err, os.ErrNotExist):
					// Defer creation to UpdatePasswordList for consistency.
					m.storeExists = false
					m.err = nil
					m.currentState = statePasswordsList

					return m, nil
				default:
					m.err = err
					m.passwordInput.SetValue("")

					return m, nil
				}
			}

			m.storeExists = true
			m.err = nil
			m.currentState = statePasswordsList

			return m, nil
		}

	case errMsg:
		m.err = msg

		return m, nil
	}

	m.passwordInput, cmd = m.passwordInput.Update(msg)

	return m, cmd
}

func EnterPasswordView(m model) string {
	wrongPassMsg := ""

	if m.err != nil {
		wrongPassMsg = fmt.Sprintf("\n\n%s\n\n", m.err)
	}

	content := fmt.Sprintf(
		"Please enter master password.%s\n\n%s\n\n%s",
		wrongPassMsg,
		keywordStyle.Render(m.passwordInput.View()),
		"(esc to quit)",
	) + "\n"
	h, v := windowStyle.GetFrameSize()

	return lipgloss.Place(
		m.vpWidth-h,
		m.vpHeight-v,
		lipgloss.Center,
		lipgloss.Center,
		formWindowStyle.Render(content),
	)
}
