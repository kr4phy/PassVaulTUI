package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	lowerLetters   = "abcdefghijklmnopqrstuvwxyz"
	upperLetters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberLetters  = "0123456789"
	specialLetters = "!@#$%^&*()-_=+[]{}\\|;:,<.>/?"
)

func setGeneratorFocus(m *model, idx int) {
	if idx < 0 {
		idx = 0
	}

	if idx > 4 {
		idx = 4
	}

	m.genFocus = idx
	m.genLengthInput.Blur()

	if idx == 0 {
		m.genLengthInput.Focus()
	}
}

func generatePassword(length int, upper, number, special bool) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than zero")
	}

	charset := []rune(lowerLetters)
	var required []rune

	if upper {
		c, err := pickRandomRune(upperLetters)
		if err != nil {
			return "", err
		}
		required = append(required, c)
		charset = append(charset, []rune(upperLetters)...)
	}

	if number {
		c, err := pickRandomRune(numberLetters)
		if err != nil {
			return "", err
		}
		required = append(required, c)
		charset = append(charset, []rune(numberLetters)...)
	}

	if special {
		c, err := pickRandomRune(specialLetters)
		if err != nil {
			return "", err
		}
		required = append(required, c)
		charset = append(charset, []rune(specialLetters)...)
	}

	if len(required) > length {
		return "", fmt.Errorf("length must be at least %d", len(required))
	}

	password := make([]rune, length)
	copy(password, required)

	for i := len(required); i < length; i++ {
		c, err := pickRandomRuneFromRunes(charset)
		if err != nil {
			return "", err
		}
		password[i] = c
	}

	for i := len(password) - 1; i > 0; i-- {
		n, err := randIndex(i + 1)
		if err != nil {
			return "", err
		}
		password[i], password[n] = password[n], password[i]
	}

	return string(password), nil
}

func pickRandomRune(chars string) (rune, error) {
	return pickRandomRuneFromRunes([]rune(chars))
}

func pickRandomRuneFromRunes(chars []rune) (rune, error) {
	if len(chars) == 0 {
		return 0, fmt.Errorf("no characters available")
	}

	idx, err := randIndex(len(chars))

	if err != nil {
		return 0, err
	}

	return chars[idx], nil
}

func randIndex(limit int) (int, error) {
	maxNum := big.NewInt(int64(limit))
	n, err := rand.Int(rand.Reader, maxNum)

	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

func UpdatePasswordGenerator(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			setGeneratorFocus(&m, 0)
			m.currentState = statePasswordsList

			return m, nil
		case "tab", "down":
			setGeneratorFocus(&m, m.genFocus+1)

			return m, nil
		case "shift+tab", "up":
			setGeneratorFocus(&m, m.genFocus-1)

			return m, nil
		case "enter", " ":
			switch m.genFocus {
			case 0:
				if msg.String() == "enter" {
					setGeneratorFocus(&m, 1)
				}
			case 1:
				m.genUseUpper = !m.genUseUpper
			case 2:
				m.genUseNumber = !m.genUseNumber
			case 3:
				m.genUseSpecial = !m.genUseSpecial
			case 4:
				length, err := strconv.Atoi(strings.TrimSpace(m.genLengthInput.Value()))

				if err != nil {
					m.err = fmt.Errorf("length must be a number")
					return m, nil
				}

				pass, genErr := generatePassword(length, m.genUseUpper, m.genUseNumber, m.genUseSpecial)

				if genErr != nil {
					m.err = genErr
					return m, nil
				}

				m.generatedPass = pass
				m.err = nil
			}

			return m, nil
		}

	case errMsg:
		m.err = msg

		return m, nil
	}

	if m.genFocus == 0 {
		var cmd tea.Cmd
		m.genLengthInput, cmd = m.genLengthInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func renderToggle(label string, enabled bool, focused bool) string {
	box := "[ ]"

	if enabled {
		box = "[x]"
	}

	line := fmt.Sprintf("%s %s", box, label)

	if focused {
		return keywordStyle.Render(line)
	}

	return line
}

func PasswordGeneratorView(m model) string {
	lenInput := m.genLengthInput.View()

	if m.genLengthInput.Focused() {
		lenInput = keywordStyle.Render(lenInput)
	}

	upperToggle := renderToggle("Uppercase", m.genUseUpper, m.genFocus == 1)
	numberToggle := renderToggle("Numbers", m.genUseNumber, m.genFocus == 2)
	specialToggle := renderToggle("Special", m.genUseSpecial, m.genFocus == 3)
	generateBtn := "[ generate ]"

	if m.genFocus == 4 {
		generateBtn = keywordStyle.Render(generateBtn)
	}

	generated := m.generatedPass

	if generated == "" {
		generated = "(not generated yet)"
	}

	errMsg := ""

	if m.err != nil {
		errMsg = fmt.Sprintf("\n\n%s", m.err)
	}

	h, v := windowStyle.GetFrameSize()
	content := "\n"
	title := titleStyle.
		Width(m.vpWidth / 2).
		AlignHorizontal(lipgloss.Center).
		Render("\nPassword Generator")
	content += fmt.Sprintf(
		textStyle.
			AlignHorizontal(lipgloss.Left).
			Render("\n\nLength %s\n\n%s\n\n%s\n\n%s\n\n%s\n\nGenerated: %s%s\n\n%s"),
		lenInput,
		upperToggle,
		numberToggle,
		specialToggle,
		generateBtn,
		keywordStyle.
			MaxWidth(m.vpWidth/2).
			Render(generated),
		errMsg,
		subtleStyle.Render("(tab/down to next, shift+tab/up back, enter/space to toggle, esc to cancel)"),
	)
	content += "\n"

	return lipgloss.Place(
		m.vpWidth-h,
		m.vpHeight-v,
		lipgloss.Center,
		lipgloss.Center,
		formWindowStyle.
			Render(title+content),
	)
}
