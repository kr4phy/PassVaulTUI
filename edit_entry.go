package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// setAddFocus updates which add-entry field is focused.
func setAddFocus(m *model, idx int) {
	if idx < 0 {
		idx = 0
	}
	if idx > 2 {
		idx = 2
	}
	m.addFocus = idx
	m.addTitleInput.Blur()
	m.addIDInput.Blur()
	m.addPassInput.Blur()
	switch idx {
	case 0:
		m.addTitleInput.Focus()
	case 1:
		m.addIDInput.Focus()
	case 2:
		m.addPassInput.Focus()
	}
}

func UpdateEditEntry(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Return to list without saving
			m.isEditing = false
			m.editIndex = -1
			m.currentState = statePasswordsList
			setAddFocus(&m, 0)
			return m, nil
		case "enter":
			if m.addFocus < 2 {
				setAddFocus(&m, m.addFocus+1)
				return m, nil
			}
			// addFocus == 2 and enter -> save
			newEntry := [3]string{m.addTitleInput.Value(), m.addIDInput.Value(), m.addPassInput.Value()}
			if m.isEditing && m.editIndex >= 0 && m.editIndex < len(m.passStorage) {
				m.passStorage[m.editIndex] = newEntry
			} else {
				m.passStorage = append(m.passStorage, newEntry)
			}
			if err := SaveToFile("data.bin", m.passStorage, deriveKey(m.masterPass)); err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.passList.SetItems(ConvertSliceToListItem(m.passStorage))
			// reset inputs for next time
			m.addTitleInput.SetValue("")
			m.addIDInput.SetValue("")
			m.addPassInput.SetValue("")
			m.isEditing = false
			m.editIndex = -1
			setAddFocus(&m, 0)
			m.currentState = statePasswordsList
			return m, nil
		case "tab", "down":
			if m.addFocus < 2 {
				setAddFocus(&m, m.addFocus+1)
				return m, nil
			}
			return m, nil
		case "shift+tab", "up":
			setAddFocus(&m, m.addFocus-1)
			return m, nil
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.addTitleInput, cmd = m.addTitleInput.Update(msg)
	cmds = append(cmds, cmd)
	m.addIDInput, cmd = m.addIDInput.Update(msg)
	cmds = append(cmds, cmd)
	m.addPassInput, cmd = m.addPassInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func EditEntryView(m model) string {
	var (
		addTitleInput string
		addIDInput    string
		addPassInput  string
	)
	if m.addTitleInput.Focused() {
		addTitleInput = keywordStyle.Render(m.addTitleInput.View())
	} else {
		addTitleInput = m.addTitleInput.View()
	}
	if m.addIDInput.Focused() {
		addIDInput = keywordStyle.Render(m.addIDInput.View())
	} else {
		addIDInput = m.addIDInput.View()
	}
	if m.addPassInput.Focused() {
		addPassInput = keywordStyle.Render(m.addPassInput.View())
	} else {
		addPassInput = m.addPassInput.View()
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
