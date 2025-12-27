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

func UpdateAddEntry(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Return to list without saving
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
			m.passStorage = append(m.passStorage, newEntry)
			if err := SaveToFile("data.bin", m.passStorage, deriveKey(m.masterPass)); err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.passList.SetItems(ConvertSliceToListItem(m.passStorage))
			// reset inputs for next time
			m.addTitleInput.SetValue("")
			m.addIDInput.SetValue("")
			m.addPassInput.SetValue("")
			setAddFocus(&m, 0)
			m.currentState = statePasswordsList
			return m, nil
		case "tab", "down":
			if m.addFocus <2 {
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

func AddEntryView(m model) string {
	content := fmt.Sprintf(
		"Add new entry\n\nTitle: %s\nID: %s\nPassword: %s\n\n(tab/down/enter to next, shift+tab/up back, enter to add item, esc to cancel)",
		m.addTitleInput.View(),
		m.addIDInput.View(),
		m.addPassInput.View(),
	)
	h, v := windowStyle.GetFrameSize()
	return lipgloss.Place(
		m.vpWidth-h,
		m.vpHeight-v,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
