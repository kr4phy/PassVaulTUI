package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var passListStyle = lipgloss.NewStyle().Padding(1, 3)

type item struct {
	title    string
	id       string
	password string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.id }
func (i item) FilterValue() string { return i.title }

func UpdatePasswordList(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	// Load data once when entering the list state.
	if m.passStorage == nil {
		if m.storeExists {
			content, err := LoadEncryptedData("data.bin", deriveKey(m.masterPass))
			if err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.passStorage = content
		} else {
			if err := SaveToFile("data.bin", [][3]string{}, deriveKey(m.masterPass)); err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.storeExists = true
			m.passStorage = make([][3]string, 0)
		}
		m.passList.SetItems(ConvertSliceToListItem(m.passStorage))
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "n":
			setAddFocus(&m, 0)
			m.addTitleInput.SetValue("")
			m.addIDInput.SetValue("")
			m.addPassInput.SetValue("")
			m.currentState = stateAddEntry
			return m, nil
		case "d", "backspace", "delete":
			idx := m.passList.Index()
			if idx >= 0 && idx < len(m.passStorage) {
				m.passStorage = append(m.passStorage[:idx], m.passStorage[idx+1:]...)
				if err := SaveToFile("data.bin", m.passStorage, deriveKey(m.masterPass)); err != nil {
					m.err = err
					return m, tea.Quit
				}
				m.passList.SetItems(ConvertSliceToListItem(m.passStorage))
			}
			return m, nil
		case "enter":
			m.chosenCredential = m.passList.SelectedItem()
			m.currentState = statePasswordDetail
		}
	case tea.WindowSizeMsg:
		h, v := passListStyle.GetFrameSize()
		m.passList.SetSize(msg.Width-h*2, msg.Height-v*2)
	}
	var cmd tea.Cmd
	m.passList, cmd = m.passList.Update(msg)
	return m, cmd
}

func PasswordListView(m model) string {
	h, v := passListStyle.GetFrameSize()
	content := passListStyle.Render(m.passList.View())
	return lipgloss.Place(
		m.vpWidth-h*2,
		m.vpHeight-v*2,
		lipgloss.Center,
		lipgloss.Bottom,
		content,
	)
}
