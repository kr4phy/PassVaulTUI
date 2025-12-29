package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func UpdateChangeMasterPass(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.changePassInput.Blur()
			m.currentState = statePasswordsList

			return m, nil
		case "enter":
			newPass := m.changePassInput.Value()

			if m.passStorage == nil {
				if m.storeExists {
					content, err := LoadEncryptedData(dataFilePath(), deriveKey(m.masterPass))
					if err != nil {
						m.err = err

						return m, nil
					}
					m.passStorage = content
				} else {
					m.passStorage = make([][3]string, 0)
				}
			}

			if err := SaveToFile(dataFilePath(), m.passStorage, deriveKey(newPass)); err != nil {
				m.err = err

				return m, nil
			}

			m.masterPass = newPass
			m.err = nil
			m.changePassInput.SetValue("")
			m.changePassInput.Blur()
			m.currentState = statePasswordsList

			return m, nil
		}
	case errMsg:
		m.err = msg

		return m, nil
	}

	m.changePassInput, cmd = m.changePassInput.Update(msg)

	return m, cmd
}

func ChangeMasterPassView(m model) string {
	errMsg := ""
	var changePassInput string

	if m.changePassInput.Focused() {
		changePassInput = keywordStyle.Render(m.changePassInput.View())
	} else {
		changePassInput = m.changePassInput.View()
	}

	if m.err != nil {
		errMsg = fmt.Sprintf("\n\n%s\n\n", m.err)
	}

	content := fmt.Sprintf(
		"Change master password.%s\n\n%s\n\n%s",
		errMsg,
		changePassInput,
		"(enter to save, esc to cancel)",
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
