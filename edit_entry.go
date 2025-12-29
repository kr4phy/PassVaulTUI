package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func setEditFocus(m *model, idx int) {
	if idx < 0 {
		idx = 0
	}

	if idx > 2 {
		idx = 2
	}

	m.editFocus = idx
	m.editTitleInput.Blur()
	m.editIDInput.Blur()
	m.editPassInput.Blur()
	switch idx {
	case 0:
		m.editTitleInput.Focus()
	case 1:
		m.editIDInput.Focus()
	case 2:
		m.editPassInput.Focus()
	}
}

func UpdateEditEntry(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.isEditing = false
			m.editIndex = -1
			m.currentState = statePasswordsList
			setEditFocus(&m, 0)

			return m, nil
		case "enter":
			if m.editFocus < 2 {
				setEditFocus(&m, m.editFocus+1)

				return m, nil
			}

			newEntry := [3]string{m.editTitleInput.Value(), m.editIDInput.Value(), m.editPassInput.Value()}

			if m.isEditing && m.editIndex >= 0 && m.editIndex < len(m.passStorage) {
				m.passStorage[m.editIndex] = newEntry
			} else {
				m.passStorage = append(m.passStorage, newEntry)
			}

			if err := SaveToFile(dataFilePath(), m.passStorage, deriveKey(m.masterPass)); err != nil {
				m.err = err

				return m, tea.Quit

			}
			m.passList.SetItems(ConvertSliceToListItem(m.passStorage))
			m.editTitleInput.SetValue("")
			m.editIDInput.SetValue("")
			m.editPassInput.SetValue("")
			m.isEditing = false
			m.editIndex = -1
			setEditFocus(&m, 0)
			m.currentState = statePasswordsList

			return m, nil
		case "tab", "down":
			if m.editFocus < 2 {
				setEditFocus(&m, m.editFocus+1)

				return m, nil
			}

			return m, nil
		case "shift+tab", "up":
			setEditFocus(&m, m.editFocus-1)

			return m, nil
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.editTitleInput, cmd = m.editTitleInput.Update(msg)
	cmds = append(cmds, cmd)
	m.editIDInput, cmd = m.editIDInput.Update(msg)
	cmds = append(cmds, cmd)
	m.editPassInput, cmd = m.editPassInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func EditEntryView(m model) string {
	var (
		addTitleInput string
		addIDInput    string
		addPassInput  string
	)

	if m.editTitleInput.Focused() {
		addTitleInput = keywordStyle.Render(m.editTitleInput.View())
	} else {
		addTitleInput = m.editTitleInput.View()
	}

	if m.editIDInput.Focused() {
		addIDInput = keywordStyle.Render(m.editIDInput.View())
	} else {
		addIDInput = m.editIDInput.View()
	}

	if m.editPassInput.Focused() {
		addPassInput = keywordStyle.Render(m.editPassInput.View())
	} else {
		addPassInput = m.editPassInput.View()
	}

	label := "Add new entry"

	if m.isEditing {
		label = "Edit entry"
	}

	content := fmt.Sprintf(
		"%s\n\nTitle %s\nID %s\nPassword %s\n\n(tab/down/enter to next, shift+tab/up back, enter to save, esc to cancel)",
		titleStyle.
			Width(m.vpWidth/2).
			Render(label),
		addTitleInput,
		addIDInput,
		addPassInput,
	)
	h, v := windowStyle.GetFrameSize()

	return lipgloss.Place(
		m.vpWidth-h,
		m.vpHeight-v,
		lipgloss.Center,
		lipgloss.Center,
		formWindowStyle.Render(content),
	)
}
