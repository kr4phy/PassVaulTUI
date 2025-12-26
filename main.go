package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	stateEnterPassword  = 0
	statePasswordsList  = 1
	statePasswordDetail = 2
)

var (
	windowStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1) // 좌우 여백
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	vpWidth          int
	vpHeight         int
	currentState     int
	passwordInput    textinput.Model
	masterPass       string
	passList         list.Model
	chosenCredential list.Item
	err              error
}

func initialModel() model {
	pInput := textinput.New()
	pInput.Placeholder = "password"
	pInput.Focus()
	pInput.CharLimit = 156
	pInput.Width = 20
	passList := ConvertSliceToListItem(LoadJSONData("data.json"))
	return model{
		vpWidth:       0,
		vpHeight:      0,
		currentState:  0,
		passwordInput: pInput,
		masterPass:    "",
		passList:      list.New(passList, list.NewDefaultDelegate(), 0, 0),
		err:           nil,
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := passListStyle.GetFrameSize()
		ih, iv := passListStyle.GetFrameSize()
		m.passList.SetSize(msg.Width-h-ih, msg.Height-v-iv)
		m.vpWidth = msg.Width
		m.vpHeight = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	switch m.currentState {
	case stateEnterPassword:
		return UpdateEnterPassword(msg, m)
	case statePasswordsList:
		return UpdatePasswordList(msg, m)
	default:
		return UpdateEnterPassword(msg, m)
	}
}

func (m model) View() string {
	var content string
	switch m.currentState {
	case stateEnterPassword:
		content = EnterPasswordView(m)
	case statePasswordsList:
		content = PasswordListView(m)
	default:
		content = EnterPasswordView(m)
	}
	return windowStyle.Render(content)
}
