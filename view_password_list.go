package main

import (
	"github.com/charmbracelet/bubbles/list"
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
		if m.passList.FilterState() == list.Filtering {
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			break
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "n":
			setAddFocus(&m, 0)
			m.addTitleInput.SetValue("")
			m.addIDInput.SetValue("")
			m.addPassInput.SetValue("")
			m.isEditing = false
			m.editIndex = -1
			m.currentState = stateAddEntry
			return m, nil
		case "e":
			idx := m.passList.Index()
			if idx >= 0 && idx < len(m.passStorage) {
				setAddFocus(&m, 0)
				m.addTitleInput.SetValue(m.passStorage[idx][0])
				m.addIDInput.SetValue(m.passStorage[idx][1])
				m.addPassInput.SetValue(m.passStorage[idx][2])
				m.isEditing = true
				m.editIndex = idx
				m.currentState = stateAddEntry
				return m, nil
			}
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
		case "s":
			m.changePassInput.SetValue("")
			m.changePassInput.Focus()
			m.err = nil
			m.currentState = stateChangeMaster
			return m, nil
		case "g":
			setGeneratorFocus(&m, 0)
			if m.genLengthInput.Value() == "" {
				m.genLengthInput.SetValue("16")
			}
			m.generatedPass = ""
			m.err = nil
			m.currentState = statePasswordGenerator
			return m, nil
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
	h, v := windowStyle.GetFrameSize()
	content := passListStyle.Render(m.passList.View())
	return lipgloss.Place(
		m.vpWidth-h,
		m.vpHeight-v,
		lipgloss.Left,
		lipgloss.Center,
		content,
	)
}
