package main

import (
	"errors"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateEnterPassword state = iota
	statePasswordsList
	statePasswordDetail
	stateEditEntry
	stateChangeMaster
	statePasswordGenerator
)

type updateHandler func(tea.Msg, model) (tea.Model, tea.Cmd)
type viewHandler func(model) string

var updateDispatch = map[state]updateHandler{
	stateEnterPassword:     UpdateEnterPassword,
	statePasswordsList:     UpdatePasswordList,
	statePasswordDetail:    UpdatePasswordDetail,
	stateEditEntry:         UpdateEditEntry,
	stateChangeMaster:      UpdateChangeMasterPass,
	statePasswordGenerator: UpdatePasswordGenerator,
}

var viewDispatch = map[state]viewHandler{
	stateEnterPassword:     EnterPasswordView,
	statePasswordsList:     PasswordListView,
	statePasswordDetail:    PasswordDetailView,
	stateEditEntry:         EditEntryView,
	stateChangeMaster:      ChangeMasterPassView,
	statePasswordGenerator: PasswordGeneratorView,
}

var (
	windowStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)
	formWindowStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)
	titleStyle   = lipgloss.NewStyle().Bold(true).AlignHorizontal(lipgloss.Center)
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	textStyle    = lipgloss.NewStyle()
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
	currentState     state
	storeExists      bool
	passwordInput    textinput.Model
	masterPass       string
	passStorage      [][3]string
	passList         list.Model
	isEditing        bool
	editIndex        int
	editTitleInput   textinput.Model
	editIDInput      textinput.Model
	editPassInput    textinput.Model
	editFocus        int
	changePassInput  textinput.Model
	genLengthInput   textinput.Model
	genUseUpper      bool
	genUseNumber     bool
	genUseSpecial    bool
	genFocus         int
	generatedPass    string
	chosenCredential list.Item
	err              error
}

func makeTextInput(placeholder string, param ...int) textinput.Model {
	inputForm := textinput.New()
	inputForm.Placeholder = placeholder

	charLimit := 156
	width := 24

	if len(param) > 0 && param[0] != 0 {
		charLimit = 156
	}

	if len(param) > 1 && param[1] != 0 {
		width = param[1]
	}

	inputForm.CharLimit = charLimit
	inputForm.Width = width

	return inputForm
}

func initialModel() model {
	pInput := textinput.New()
	pInput.Placeholder = "Enter master password"
	pInput.Focus()
	pInput.CharLimit = 156
	pInput.Width = 24
	pInput.EchoMode = textinput.EchoPassword
	pInput.EchoCharacter = '*'
	passList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	passList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new entry")),
			key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit entry")),
			key.NewBinding(key.WithKeys("d", "backspace", "delete"), key.WithHelp("d/backspace/delete", "remove")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open detail")),
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "change master")),
			key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "generate password")),
		}
	}

	var storeExists bool

	if _, err := os.Stat(dataFilePath()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			storeExists = false
		} else {
			log.Println("Uncaught os.Stat() error:", err)
			storeExists = false
		}
	} else {
		storeExists = true
	}

	editTitle := makeTextInput("title")
	editID := makeTextInput("id")
	editPass := makeTextInput("password")
	editPass.EchoMode = textinput.EchoPassword
	editPass.EchoCharacter = '*'
	changePass := makeTextInput("Enter new master password")
	changePass.EchoMode = textinput.EchoPassword
	changePass.EchoCharacter = '*'
	genLength := makeTextInput("length", 4, 6)
	genLength.SetValue("16")

	return model{
		vpWidth:          0,
		vpHeight:         0,
		currentState:     stateEnterPassword,
		storeExists:      storeExists,
		passwordInput:    pInput,
		masterPass:       "",
		passStorage:      nil,
		passList:         passList,
		isEditing:        false,
		editIndex:        -1,
		editTitleInput:   editTitle,
		editIDInput:      editID,
		editPassInput:    editPass,
		editFocus:        0,
		changePassInput:  changePass,
		genLengthInput:   genLength,
		genUseUpper:      true,
		genUseNumber:     true,
		genUseSpecial:    true,
		genFocus:         0,
		generatedPass:    "",
		chosenCredential: item{},
		err:              nil,
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := passListStyle.GetFrameSize()
		m.passList.SetSize(msg.Width-h*2, msg.Height-v*2)
		m.vpWidth = msg.Width
		m.vpHeight = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if h, ok := updateDispatch[m.currentState]; ok {
		return h(msg, m)
	}

	return m, nil
}

func (m model) View() string {
	var content string

	if h, ok := viewDispatch[m.currentState]; ok {
		content = h(m)
	}

	return windowStyle.Render(content)
}
